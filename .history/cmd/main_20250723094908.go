package main

import "time"

func init() {
	//load config info from file env.dev.yml or env.pro.yml
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return
	}
	time.Local = location
}

func main() {

}
