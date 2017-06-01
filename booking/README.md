Hotel Booking Example
===============================

The Hotel Booking example app demonstrates ([browse the source](https://github.com/revel/examples/tree/master/booking)):

* Using an SQL (SQLite) database and configuring the Revel DB module.
* Using the third party [GORP](https://github.com/coopernurse/gorp) *ORM-ish* library
* [Interceptors](../manual/interceptors.html) for checking that an user is logged in.
* Using [validation](../manual/validation) and displaying inline errors


Here's a quick summary of the structure
```
	booking/app/
		models		   # Structs and validation.
			booking.go
			hotel.go
			user.go

		controllers
			init.go    # Register all of the interceptors.
			gorp.go    # A plugin for setting up Gorp, creating tables, and managing transactions.
			app.go     # "Login" and "Register new user" pages
			hotels.go  # Hotel searching and booking

		views
			...
```

# Database Install and Setup
This example used [sqlite](https://www.sqlite.org/), (Alternatively can use mysql, postgres, etc.)

## sqlite Installation

- The booking app uses [go-sqlite3](https://github.com/mattn/go-sqlite3) database driver, which depends on the C library

### Install sqlite on OSX:

1. Install [Homebrew](http://mxcl.github.com/homebrew/) if you don't already have it.
2. Install pkg-config and sqlite3:

~~~
$ brew install pkgconfig sqlite3
~~~

### Install sqlite on Ubuntu:
```sh
$ sudo apt-get install sqlite3 libsqlite3-dev
```

Once SQLite is installed, it will be possible to run the booking app:
```sh
	$ revel run github.com/revel/examples/booking
```

## Database / Gorp Plugin

[`app/controllers/gorp.go`](https://github.com/revel/examples/blob/master/booking/app/controllers/gorp.go) defines `GorpPlugin`, which is a plugin that does a couple things:

* **`OnAppStart`** -  Uses the DB module to open a SQLite in-memory database, create the `User`, `Booking`, and `Hotel` tables, and insert some test records.
* **BeforeRequest** -  Begins a transaction and stores the Transaction on the Controller
* **AfterRequest** -  Commits the transaction, or [panics](https://github.com/golang/go/wiki/PanicAndRecover) if there was an error.
* **OnException** -  Rolls back the transaction


## Interceptors

[`app/controllers/init.go`](https://github.com/revel/examples/blob/master/booking/app/controllers/init.go) 
registers the [interceptors](../manual/interceptors.html) that runs before each action:

```go
func init() {
	revel.OnAppStart(Init)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	revel.InterceptMethod(Hotels.checkUser, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}
```

As an example, `checkUser` looks up the username in the `session` and `redirect`s
the user to log in if they do not have a `session` cookie.

```go
func (c Hotels) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(Application.Index)
	}
	return nil
}
```

[Check out the user management code in app.go](https://github.com/revel/examples/blob/master/booking/app/controllers/app.go)

## Validation

The booking app does quite a bit of validation.

For example, here is the routine to validate a booking, from
[models/booking.go](https://github.com/revel/examples/blob/master/booking/app/models/booking.go):

```go
func (booking Booking) Validate(v *revel.Validation) {
	v.Required(booking.User)
	v.Required(booking.Hotel)
	v.Required(booking.CheckInDate)
	v.Required(booking.CheckOutDate)

	v.Match(b.CardNumber, regexp.MustCompile(`\d{16}`)).
		Message("Credit card number must be numeric and 16 digits")

	v.Check(booking.NameOnCard,
		revel.Required{},
		revel.MinSize{3},
		revel.MaxSize{70},
	)
}
```

Revel applies the validation and records errors using the name of the
validated variable (unless overridden).  For example, `booking.CheckInDate` is
required; if it evaluates to the zero date, Revel stores a `ValidationError` in
the validation context under the key "booking.CheckInDate".

Subsequently, the
[Hotels/Book.html](https://github.com/revel/examples/blob/master/booking/app/views/Hotels/Book.html)
template can access them using the [`field`](../manual/templates.html#field) helper:

```
{% capture ex %}{% raw %}
{{with $field := field "booking.CheckInDate" .}}
<p class="{{$field.ErrorClass}}">
    <strong>Check In Date:</strong>
    <input type="text" size="10" name="{{$field.Name}}" class="datepicker" value="{{$field.Flash}}">
    * <span class="error">{{$field.Error}}</span>
ss</p>
{{end}}
{% endraw %}{% endcapture %}
{% highlight htmldjango %}{{ex}}{% endhighlight %} 
```


The [`field`](../manual/templates.html#field) template helper looks for errors in the validation context, using
the field name as the key.
