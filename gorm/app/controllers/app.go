package controllers

import (
	"github.com/revel/revel"
	"modules/gorm/app"
	"testgorm/app/models"
)

type App struct {
	*revel.Controller
	gorm.GormController
	// Transactional
	// gorm.GormTransactionController
}

func (c App) Index() revel.Result {
	var users = []models.User{}
	c.Txn.Find(&users)
	return c.RenderJSON(users)
}
