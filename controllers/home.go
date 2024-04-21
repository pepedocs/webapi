package controllers

import (
	"fmt"
	"net/http"
)

type IHomeService interface {
	Greet() string
}

type HomeController struct {
	HomeSvc IHomeService
}

func (c *HomeController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, c.HomeSvc.Greet())
}
