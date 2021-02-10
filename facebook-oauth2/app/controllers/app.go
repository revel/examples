package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/revel/examples/facebook-oauth2/app/models"
	"github.com/revel/revel"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type Application struct {
	*revel.Controller
}

// The following keys correspond to a test application
// registered on Facebook, and associated with the loisant.org domain.
// You need to bind loisant.org to your machine with /etc/hosts to
// test the application locally.

var Facebook = &oauth2.Config{ //nolint:gochecknoglobals
	ClientID:     "943076975742162",
	ClientSecret: "d3229ebe3501771344bb0f2db2324014",
	Scopes:       []string{},
	Endpoint:     facebook.Endpoint,
	RedirectURL:  "http://loisant.org:9000/Application/Auth",
}

func (c Application) Index() revel.Result {
	u := c.connected()
	me := map[string]interface{}{}
	if u != nil && u.AccessToken != "" {
		resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
			url.QueryEscape(u.AccessToken))
		if err != nil {
			c.Log.Error("Failed HTTP GET", "error", err)

			return nil
		}

		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
			c.Log.Error("json decode error", "error", err)
		}

		c.Log.Info("Data fetched", "data", me)
	}

	authURL := Facebook.AuthCodeURL("state", oauth2.AccessTypeOffline)

	return c.Render(me, authURL)
}

func (c Application) Auth(code string) revel.Result {
	tok, err := Facebook.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.Log.Error("Exchange error", "error", err)
		return c.Redirect(Application.Index)
	}

	user := c.connected()
	user.AccessToken = tok.AccessToken

	return c.Redirect(Application.Index)
}

func setuser(c *revel.Controller) revel.Result {
	var user *models.User
	if _, ok := c.Session["uid"]; ok {
		uid, _ := strconv.ParseInt(c.Session["uid"].(string), 10, 0)
		user = models.GetUser(int(uid))
	}

	if user == nil {
		user = models.NewUser()
		c.Session["uid"] = fmt.Sprintf("%d", user.UID)
	}
	c.ViewArgs["user"] = user

	return nil
}

func init() {
	revel.InterceptFunc(setuser, revel.BEFORE, &Application{})
}

func (c Application) connected() *models.User {
	return c.ViewArgs["user"].(*models.User)
}
