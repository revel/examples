Chat App Demo
=========================
The `Chat` app demonstrates ([browse the source](https://github.com/revel/samples/tree/master/chat)):

* Using channels to implement a chat room with a [publish-subscribe](http://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern) model.
* Using both `Comet` and [Websockets](../manual/websockets.html)

Here's a quick summary of the structure:
```
	chat/app/
		chatroom	       # Chat room routines
			chatroom.go

		controllers
			app.go         # The login screen, allowing user to choose from supported technologies
			refresh.go     # Handlers for the "Active Refresh" chat demo
			longpolling.go # Handlers for the "Long polling" ("Comet") chat demo
			websocket.go   # Handlers for the "Websocket" chat demo

		views
			                # HTML and Javascript

```
# Chat Room Background

First, let's look at how the chat room is implemented, in
[`app/chatroom/chatroom.go`](https://github.com/revel/samples/blob/master/chat/app/chatroom/chatroom.go).

The chat room runs as an independent `go-routine`, started on initialization:

```go
func init() {
	go chatroom()
}
```

The `chatroom()` function simply selects on three channels to execute the requested action.

```go
var (
	// Send a channel here to get room events back.  It will send the entire
	// archive initially, and then new messages as they come in.
	subscribe = make(chan (chan<- Subscription), 10)
	// Send a channel here to unsubscribe.
	unsubscribe = make(chan (<-chan Event), 10)
	// Send events here to publish them.
	publish = make(chan Event, 10)
)

func chatroom() {
	archive := list.New()
	subscribers := list.New()

	for {
		select {
		case ch := <-subscribe:
			// Add subscriber to list and send back subscriber channel + chat log.
		case event := <-publish:
			// Send event to all subscribers and add to chat log.
		case unsub := <-unsubscribe:
			// Remove subscriber from subscriber list.
		}
	}
}
```

Let's examine how each of those channel functions are implemented.

### Subscribe 

```go
case ch := <-subscribe:
    var events []Event
    for e := archive.Front(); e != nil; e = e.Next() {
        events = append(events, e.Value.(Event))
    }
    subscriber := make(chan Event, 10)
    subscribers.PushBack(subscriber)
    ch <- Subscription{events, subscriber}
```

A `Subscription` is created with two properties:

* The chat log (archive)
* A channel that the subscriber can listen on to get new messages.

The `Subscription` is then sent to the channel that subscriber supplied.


### Publish

```go
case event := <-publish:
    for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
        ch.Value.(chan Event) <- event
    }
    if archive.Len() >= archiveSize {
        archive.Remove(archive.Front())
    }
    archive.PushBack(event)
```

The `Published event` is sent to the subscribers' channels one by one.  
- The `event` is added to the `archive`, which is trimmed if necessary.

### Unsubscribe

```go
case unsub := <-unsubscribe:
    for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
        if ch.Value.(chan Event) == unsub {
            subscribers.Remove(ch)
        }
    }
```

The `Subscriber` channel is removed from the list.

## Handlers

Now that the `Chat Room` channels exist, lets examine how the handlers
expose that functionality using different techniques.

### Active Refresh

The Active Refresh chat room javascript refreshes the page every five seconds to
get any new messages (see [`Refresh/Room.html`](https://github.com/revel/samples/blob/master/chat/app/views/Refresh/Room.html)):

```js
// Scroll the messages panel to the end
var scrollDown = function() {
    $('#thread').scrollTo('max');
}

// Reload the whole messages panel
var refresh = function() {
    $('#thread').load('/refresh/room?user={{.user}} #thread .message', function() {
        scrollDown();
    })
}

// Call refresh every 5 seconds
setInterval(refresh, 5000);
```

This is the handler to serve the above in [`app/controllers/refresh.go`](https://github.com/revel/revel/tree/master/samples/chat/app/controllers/refresh.go):

```go
func (c Refresh) Room(user string) revel.Result {
	subscription := chatroom.Subscribe()
	defer subscription.Cancel()
	events := subscription.Archive
	for i, _ := range events {
		if events[i].User == user {
			events[i].User = "you"
		}
	}
	return c.Render(user, events)
}
```


It subscribes to the chatroom and passes the archive to the template to be
rendered (after changing the user name to "you" as necessary).



### Long Polling with Comet

The Long Polling chat room (see [`LongPolling/Room.html`](https://github.com/revel/revel/tree/master/samples/chat/app/views/LongPolling/Room.html))
makes an ajax request that the server keeps open until a new message comes in. The javascript uses a
`lastReceived` timestamp to tell the server the last message it knows about.

``` js 
var lastReceived = 0;
var waitMessages = '/longpolling/room/messages?lastReceived=';
var say = '/longpolling/room/messages?user={{.user}}';

$('#send').click(function(e) {
    var message = $('#message').val();
    $('#message').val('');
    $.post(say, {message: message});
});

// Retrieve new messages
var getMessages = function() {
    $.ajax({
        url: waitMessages + lastReceived,
        success: function(events) {
            $(events).each(function() {
                display(this);
                lastReceived = this.Timestamp;
            });
            getMessages();
        },
        dataType: 'json'
    });
}
getMessages();
```

The handler for the above in [`app/controllers/longpolling.go`](https://github.com/revel/samples/blob/master/chat/app/controllers/longpolling.go)

```go
func (c LongPolling) WaitMessages(lastReceived int) revel.Result {
	subscription := chatroom.Subscribe()
	defer subscription.Cancel()

	// See if anything is new in the archive.
	var events []chatroom.Event
	for _, event := range subscription.Archive {
		if event.Timestamp > lastReceived {
			events = append(events, event)
		}
	}

	// If we found one, grand.
	if len(events) > 0 {
		return c.RenderJson(events)
	}

	// Else, wait for something new.
	event := <-subscription.New
	return c.RenderJson([]chatroom.Event{event})
}
```


In this implementation, it can simply block on the subscription channel, 
assuming it has already sent back everything in the archive.

### Websocket

The Websocket chat room (see  [WebSocket/Room.html](https://github.com/revel/samples/blob/master/chat/app/views/WebSocket/Room.html#L51))
opens a [websocket](../manual/websockets.html) connection as soon as the
user has loaded the page.

```js
// Create a socket
var socket = new WebSocket('ws://127.0.0.1:9000/websocket/room/socket?user={{.user}}');

// Message received on the socket
socket.onmessage = function(event) {
    display(JSON.parse(event.data));
}

$('#send').click(function(e) {
    var message = $('#message').val();
    $('#message').val('');
    socket.send(message);
});
```

The first thing to do is to subscribe to new events, join the room, and send
down the archive.  Here is what [websocket.go](https://github.com/revel/revel/tree/master/samples/chat/app/controllers/websocket.go#L17) looks like:

```go
func (c WebSocket) RoomSocket(user string, ws *websocket.Conn) revel.Result {
	// Join the room.
	subscription := chatroom.Subscribe()
	defer subscription.Cancel()

	chatroom.Join(user)
	defer chatroom.Leave(user)

	// Send down the archive.
	for _, event := range subscription.Archive {
		if websocket.JSON.Send(ws, &event) != nil {
			// They disconnected
			return nil
		}
	}
	....
```


Next, we have to listen for new events from the subscription.  However, the
websocket library only provides a blocking call to get a new frame.  To select
between them, we have to wrap it ([websocket.go](https://github.com/revel/revel/tree/master/samples/chat/app/controllers/websocket.go#L33)):

```go
// In order to select between websocket messages and subscription events, we
// need to stuff websocket events into a channel.
newMessages := make(chan string)
go func() {
    var msg string
    for {
        err := websocket.Message.Receive(ws, &msg)
        if err != nil {
            close(newMessages)
            return
        }
        newMessages <- msg
    }
}()
```

Now we can select for new websocket messages on the `newMessages` channel.

The last bit does exactly that -- it waits for a new message from the websocket
(if the user has said something) or from the subscription (someone else in the
chat room has said something) and propagates the message to the other.

```go
// Now listen for new events from either the websocket or the chatroom.
for {
    select {
    case event := <-subscription.New:
        if websocket.JSON.Send(ws, &event) != nil {
            // They disconnected.
            return nil
        }
    case msg, ok := <-newMessages:
        // If the channel is closed, they disconnected.
        if !ok {
            return nil
        }

        // Otherwise, say something.
        chatroom.Say(user, msg)
    }
}
return nil

```

> [websocket.go](https://github.com/revel/revel/tree/master/samples/chat/app/controllers/websocket.go#L48)

If we detect the websocket channel has closed, then we just return nil.

