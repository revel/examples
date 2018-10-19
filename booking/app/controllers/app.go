package controllers

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/revel/revel"

	"database/sql"
	"github.com/revel/examples/booking/app/models"
	"github.com/revel/examples/booking/app/routes"
	"github.com/revel/modules/orm/gorp/app/controllers"
)

type Application struct {
	gorpController.Controller
}

func (c Application) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

func (c Application) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username.(string))
	}
	
	return nil
}

func (c Application) getUser(username string) (user *models.User) {
	user = &models.User{}
	_,  err := c.Session.GetInto("fulluser", user, false)
	if user.Username == username {
		return user
	}

	err = c.Txn.SelectOne(user, c.Db.SqlStatementBuilder.Select("*").From("User").Where("Username=?", username))
	if err != nil {
		if err != sql.ErrNoRows {
			//c.Txn.Select(user, c.Db.SqlStatementBuilder.Select("*").From("User").Limit(1))
			count, _ := c.Txn.SelectInt(c.Db.SqlStatementBuilder.Select("count(*)").From("User"))
			c.Log.Error("Failed to find user", "user", username, "error",err, "count", count)
		}
		return nil
	}
	c.Session["fulluser"] = user
	return
}

func (c Application) Index() revel.Result {
	if c.connected() != nil {
		return c.Redirect(routes.Hotels.Index())
	}
	c.Flash.Error("Please log in first")
	return c.Render()
}

func (c Application) Register() revel.Result {
	return c.Render()
}

func (c Application) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		MessageKey("Password does not match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Application.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	err := c.Txn.Insert(&user)
	if err != nil {
		panic(err)
	}

	c.Session["user"] = user.Username
	c.Flash.Success("Welcome, " + user.Name)
	return c.Redirect(routes.Hotels.Index())
}

func (c Application) Login(username, password string, remember bool) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(routes.Hotels.Index())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.Application.Index())
}

func (c Application) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Application.Index())
}
