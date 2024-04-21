package controllers

import (
	"fmt"
	"net/http"
)

type ISwatchTimeService interface {
	GetInternetTime() (string, error)
}

type SwatchTimeController struct {
	SwatchTimeSvc ISwatchTimeService
}

func (c *SwatchTimeController) GetInternetTime(w http.ResponseWriter, r *http.Request) {
	swatchTime, _ := c.SwatchTimeSvc.GetInternetTime()
	// Todo: Use views here as it becomes necessary
	fmt.Fprint(w, swatchTime)
}
