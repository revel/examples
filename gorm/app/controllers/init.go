package controllers

import (
	"github.com/revel/revel"
	"modules/gorm/app"
	"testgorm/app/models"
)

func InitializeDB() {
	gorm.InitDB()
	gorm.DB.AutoMigrate(&models.User{})
	var firstUser = models.User{Name: "Demo", Email: "demo@demo.com"}
	firstUser.SetNewPassword("demo")
	firstUser.Active = true
	gorm.DB.Create(&firstUser)
}

func init() {
	revel.OnAppStart(func() {

		InitializeDB()
		revel.InterceptMethod((*gorm.GormController).Begin, revel.BEFORE)
		revel.InterceptMethod((*gorm.GormController).Commit, revel.AFTER)
		revel.InterceptMethod((*gorm.GormController).Rollback, revel.FINALLY)

		// Transaccional
		// revel.InterceptMethod((*gorm.GormTransactionController).Begin, revel.BEFORE)
		// revel.InterceptMethod((*gorm.GormTransactionController).Commit, revel.AFTER)
		// revel.InterceptMethod((*gorm.GormTransactionController).Rollback, revel.FINALLY)
	})
}
