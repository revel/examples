//go:generate swagger generate spec -o swagger.json

// Package classification Swagger Example.
// Swagger Example
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
package controllers

import (
	"github.com/revel/revel"
)

type Application struct {
	*revel.Controller
}

// swagger:route GET / Index
//
// Index
//
//
//
//     Produces:
//     - text/html
//
//     Schemes: https, http, ws
//
//
//     Responses:
//       200: Success

// swagger:operation POST /some/route route description
//
// Route Description
//
//
// ---
// produces:
// - text/html
// responses:
//   '200':
//     description: Success
func (c Application) Index() revel.Result {

	return c.Render()
}

func (c Application) Destroy() {
	c.Controller.Destroy()
}

// swagger:route GET /demo enter demo
//
// Enter Demo
//
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - text/html
//
//     Schemes: https
//
//
//     Responses:
//       200: Success
//       401: Invalid User

// swagger:operation GET /demo demo
//
// Enter Demo
//
//
// ---
// produces:
// - text/html
// parameters:
// - name: user
//   in: formData
//   description: user
//   required: true
//   type: string
// - name: demo
//   in: formData
//   description: demo
//   required: true
//   type: string
// responses:
//   '200':
//     description: Success
//   '401':
//     description: Invalid User
func (c Application) EnterDemo(user, demo string) revel.Result {
	c.Validation.Required(user)
	c.Validation.Required(demo)

	if c.Validation.HasErrors() {
		c.Flash.Error("Please choose a nick name and the demonstration type.")
		return c.Redirect(Application.Index)
	}

	switch demo {
	case "refresh":
		return c.Redirect("/refresh?user=%s", user)
	case "longpolling":
		return c.Redirect("/longpolling/room?user=%s", user)
	case "websocket":
		return c.Redirect("/websocket/room?user=%s", user)
	}
	return nil
}
