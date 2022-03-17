package main

import (
	"awesomeProject2/database"
	"awesomeProject2/router"
	"awesomeProject2/setting"
	"awesomeProject2/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {

	for {
		run()

		time.Sleep(time.Second * 2)
	}

}

func catchPanic() {
	message := recover()
	if message != nil {
		utils.Logger.Println(message)
		fmt.Printf("%v\n", message)
	}

	fmt.Println(time.Now().String()[:19] + ": SERVER MUST RELOAD")
}

func run() {

	defer catchPanic()

	// Загрузка файла конфигурации
	option := setting.Load("setting.json")

	// Подключиться к БД
	database.Connect(option)

	_ = router.Initialized().Run(option.Address + ":" + option.Port)

}

type M struct {
	ID       int32
	Name     string
	Link     string
	Children []M
}

func index2(c *gin.Context) {

	type requestData struct {
		Date string `json:"Date"`
	}

	var data requestData

	e := c.BindJSON(&data)
	if e != nil {
		fmt.Println(e)
		c.JSON(400, nil)
		return
	}

	fmt.Println(data)

	c.JSON(200, gin.H{
		"Users":   "HELLO",
		"IsAdmin": true,
		"Date":    "2022-02-18",
	})
}

func upload(c *gin.Context) {

	form, e := c.MultipartForm()
	if e != nil {
		utils.Logger.Println(e)
		c.JSON(400, nil)
		return
	}

	files := form.File["MyFiles"]

	for _, file := range files {
		e = c.SaveUploadedFile(file, "files/"+file.Filename)
		if e != nil {
			utils.Logger.Println(e)
		}
	}

	c.JSON(200, nil)
}
