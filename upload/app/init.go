package app

import (
	"github.com/revel/revel"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalid    Error = "invalid dict call"
	ErrStringKeys Error = "dict keys must be strings"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)

	// Helper for calling a template with several pipeline parameters.
	// Example input: key1 value1 key2 value2.
	// Example use: {{template "button.html" dict "dot" . "class" "active"}}
	revel.TemplateFuncs["dict"] = func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, ErrInvalid
		}

		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			if key, ok := values[i].(string); ok {
				dict[key] = values[i+1]
			} else {
				return nil, ErrStringKeys
			}
		}

		return dict, nil
	}
}

// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not.
func HeaderFilter(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
