package main

import (
	"sound_qr_services/delivery/http/routes"
	"time"

	"github.com/labstack/echo/v4"
)

func init() {
	//load config info from file env.dev.yml or env.pro.yml
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return
	}
	time.Local = location
}

func main() {
	echo := echo.New()
	routes.NewAPI(echo)
}
