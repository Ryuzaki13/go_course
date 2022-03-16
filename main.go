package main

import (
	"awesomeProject2/db"
	"awesomeProject2/setting"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitLogger()

	opt := setting.Load("setting.json")
	e := db.Connect(opt)
	if e != nil {
		fmt.Println(e)
		return
	}

	router := gin.Default()
	router.LoadHTMLGlob("template/*")
	router.Static("assets", "assets")
	router.GET("/", index)
	router.POST("/user", index2)

	router.POST("/api/product", getProducts)
	router.DELETE("/api/product", delProduct)

	_ = router.Run(opt.Address + ":" + opt.Port)
}

type M struct {
	ID       int32
	Name     string
	Link     string
	Children []M
}

func index(c *gin.Context) {
	user := db.User{}

	users := user.SelectAll()

	menu := make([]M, 0)
	list := db.SelectMenu()

	for _, item := range list {
		if item.Parent == 0 {
			for i := range item.Name {
				menu = append(menu, M{
					//ID: item.ID[i],
					Name: item.Name[i],
					Link: item.Link[i],
				})
			}
		} else {
			for _, a := range menu {
				if a.ID == item.Parent {
					for i := range item.Name {
						a.Children = append(a.Children, M{
							//ID: item.ID[i],
							Name: item.Name[i],
							Link: item.Link[i],
						})
					}
				}
			}
		}
	}

	c.HTML(200, "index", gin.H{
		"Users":   users,
		"Title":   "Сайтик",
		"IsAdmin": false,
		"Menu":    menu,
	})
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
		db.Logger.Println(e)
		c.JSON(400, nil)
		return
	}

	files := form.File["MyFiles"]

	for _, file := range files {
		e = c.SaveUploadedFile(file, "files/"+file.Filename)
		if e != nil {
			db.Logger.Println(e)
		}
	}

	c.JSON(200, nil)
}

func getProducts(c *gin.Context) {
	type input struct {
		Search string `json:"Search"`
	}

	i := input{}
	e := c.BindJSON(&i)
	if e != nil {
		db.Logger.Println(e)
		c.JSON(400, nil)
		return
	}

	c.JSON(200, db.SearchProduct(i.Search))
}

func delProduct(c *gin.Context) {
	type input struct {
		ID int `json:"id"`
	}

	i := input{}
	e := c.BindJSON(&i)
	if e != nil {
		db.Logger.Println(e)
		c.JSON(400, nil)
		return
	}

	e = db.DeleteProduct(i.ID)
	if e != nil {
		c.JSON(400, nil)
		return
	}

	c.JSON(200, nil)
}
