package controllers

import "github.com/revel/revel"

type Application struct {
	*revel.Controller
}

func (c Application) Index() revel.Result {
	// Localization information
	c.ViewArgs["acceptLanguageHeader"] = c.Request.Header.Get("Accept-Language")
	c.ViewArgs["acceptLanguageHeaderParsed"] = c.Request.AcceptLanguages.String()
	c.ViewArgs["acceptLanguageHeaderMostQualified"] = c.Request.AcceptLanguages[0]
	c.ViewArgs["controllerCurrentLocale"] = c.Request.Locale

	// Controller-resolves messages
	c.ViewArgs["controllerGreeting"] = c.Message("greeting")
	c.ViewArgs["controllerGreetingName"] = c.Message("greeting.name")
	c.ViewArgs["controllerGreetingSuffix"] = c.Message("greeting.suffix")
	c.ViewArgs["controllerGreetingFull"] = c.Message("greeting.full")
	c.ViewArgs["controllerGreetingWithArgument"] = c.Message("greeting.full.name", "Steve Buscemi")

	return c.Render()
}
