Swagger Demo
=========================
The `Swagger` app demonstrates ([browse the source](https://github.com/revel/samples/tree/master/swagger)):

* Using [go-swagger](https://github.com/go-swagger/go-swagger) library to generate a spec based on the [chat](https://github.com/revel/samples/tree/master/chat) example

Here's a quick summary of the structure:
```
	swagger/app/
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
# Swagger Meta

#### Inside of your app/controller/app.go file, at the top put:
*The comment lines are necessary*
```
//go:generate swagger generate spec -o swagger.json
```

#### Make sure you put a blank new line and then add this:
```
// Package classification Some Example API.
// Example API
//
//
//
//     Schemes: https
//     Host: api.somedomain.com
//     BasePath: /
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Name<email@somehwere.com> https://www.somewhere.com
//
//     Consumes:
//     - application/json
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - text/html
//
//
//
//
// swagger:meta
```

# Swagger Route

#### Inside each of your controllers and above each route add the following:

```
// swagger:route POST /some/route route description
//
// Route description
//
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - text/html
//
//     Schemes: https, http, ws
//
//
//     Responses:
//       200: Success
//       401: Invalid Info

// swagger:operation POST /some/route route description
//
// Route Description
//
//
// ---
// produces:
// - text/html
// parameters:
// - name: some_param
//   in: formData
//   description: example param
//   required: true
//   type: string
// responses:
//   '200':
//     description: Success
//   '401':
//     description: Invalid Info
```

# Generating

#### You will need to cd into the app/controllers directory and then run:
```
go generate
```

This will put a swagger.json file inside of your app/controllers folder. Feel free to add that file to ignore for git and have CI run the go generate on compile and you can now automate Swagger inside of your API.
