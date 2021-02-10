package controllers

import "github.com/revel/revel"

func init() { //nolint:gochecknoinits
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	revel.InterceptMethod(Hotels.checkUser, revel.BEFORE)
}
