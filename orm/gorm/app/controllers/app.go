package controllers

import (
	"github.com/revel/examples/orm/gorm/app/models"
	gormc "github.com/revel/modules/orm/gorm/app/controllers"
	"github.com/revel/revel"
)

type App struct {
	gormc.TxnController
}

func (c App) Index() revel.Result {
	users := []models.User{}
	c.Txn.Find(&users)
	return c.RenderJSON(users)
}
