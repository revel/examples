package controllers

import (
	"github.com/revel/examples/gorm/app/models"

	"github.com/revel/modules/gorm/app"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
	gorm.GormController
}

func (c App) Index() revel.Result {
	var users = []models.User{}
	c.Txn.Find(&users)
	return c.RenderJSON(users)
}
