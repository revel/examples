Twitter OAuth with Revel example
==================================

The `twitter-oauth` app uses the `mrjones/oauth` library to demonstrate:

* How to do the `OAuth` to authenticate your app to use a `Twitter` account.
* Fetching mentions for that Twitter account.
* Tweeting on behalf of that Twitter account.

The core contents of the app:
```sh
	twitter-oauth/app/
		models
			user.go   # User struct and in-memory data store
		controllers
			app.go    # All code
```
[Browse the code on Github](https://github.com/revel/samples/tree/master/twitter-oauth)

## OAuth Overview

The `OAuth` process is governed by this configuration:

```go
var TWITTER = oauth.NewConsumer(
	"VgRrevelRunFooBarTruingindreamzw",
	"l8lOLyIF3peCFEvrEoTc8h4oFwieAFgPM6eegibberish",
	oauth.ServiceProvider{
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	},
)
```

An overview of the process:

1. The app generates a "request token" and sends the user to Twitter.
2. The user authorizes the app.
3. Twitter redirects the user to the provided redirect url, including an
   "verifier" in the parameters.
4. The app constructs a request to Twitter using the "request token" and
   the "verifier", to which Twitter returns the "access token".
5. The app henceforth uses the access token to operate Twitter on the user's behalf.


