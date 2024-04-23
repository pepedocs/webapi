package controllers

import (
	"fmt"
	"net/http"
)

type iHomeService interface {
	Greet() string
}

type HomeController struct {
	HomeSvc iHomeService
}

func (c *HomeController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, c.HomeSvc.Greet())
}
