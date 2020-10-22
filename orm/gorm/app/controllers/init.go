package controllers

import (
	"github.com/revel/examples/orm/gorm/app/models"
	gorm "github.com/revel/modules/orm/gorm/app"
	"github.com/revel/revel"
)

func initializeDB() {
	gorm.DB.AutoMigrate(&models.User{})
	firstUser := models.User{Name: "Demo", Email: "demo@demo.com"}
	firstUser.SetNewPassword("demo")
	firstUser.Active = true
	gorm.DB.Create(&firstUser)
}

func init() {
	revel.OnAppStart(initializeDB)
}
