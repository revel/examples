package controllers

import (
	"github.com/revel/revel"
	gormc "github.com/revel/modules/orm/gorm/app/controllers"
	"github.com/revel/examples/orm/gorm/app/models"
)

type App struct {
	gormc.TxnController
}

func (c App) Index() revel.Result {
	var users = []models.User{}
	c.Txn.Find(&users)
	return c.RenderJSON(users)
}
