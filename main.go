package main

import (
	"awesomeProject2/db"
	"awesomeProject2/setting"
	"awesomeProject2/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
)

var sessionMap map[string]db.Session

var adminJS template.JS

func main() {
	db.InitLogger()

	opt := setting.Load("setting.json")
	e := db.Connect(opt)
	if e != nil {
		fmt.Println(e)
		return
	}

	adminJS = template.JS(utils.LoadAssets("secret/admin.js"))

	sessionMap = make(map[string]db.Session)

	db.LoadSession(sessionMap)

	word := sessions.NewCookieStore([]byte("SecretX"))

	router := gin.Default()

	router.Use(sessions.Sessions("hello", word))

	router.LoadHTMLGlob("template/*")
	router.Static("assets", "assets")
	router.GET("/", index)
	router.POST("/user", index2)

	router.POST("/api/product", getProducts)
	router.DELETE("/api/product", delProduct)

	router.POST("/reg", reg)
	router.POST("/login", login)
	router.POST("/logout", logout)

	_ = router.Run(opt.Address + ":" + opt.Port)
}

func reg(c *gin.Context) {
	var user db.User
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(403, nil)
		return
	}

	user.Password, e = utils.Encrypt(user.Password)
	if e != nil {
		c.JSON(400, nil)
		return
	}

	if user.Insert() {
		c.JSON(200, nil)
		return
	}

	c.JSON(400, nil)
}

func login(c *gin.Context) {

	session := sessions.Default(c)

	var user db.User
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(403, nil)
		return
	}

	user.Password, e = utils.Encrypt(user.Password)
	if e != nil {
		c.JSON(400, nil)
		return
	}

	if user.LogIn() {
		hash, ok := db.CreateSession(&user, sessionMap)
		if ok {
			session.Set("SessionSecretKey", hash)
			e = session.Save()
			if e != nil {
				db.Logger.Println(e)
			}

			c.JSON(200, nil)

			return
		}
	}

	c.JSON(400, nil)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)

	_, ok := session.Get("SessionSecretKey").(string)
	if ok {
		session.Clear()
		_ = session.Save()
		c.SetCookie("hello", "", -1, "/", c.Request.URL.Hostname(), false, true)
		session.Delete("SessionSecretKey")
	}

	c.Redirect(301, c.Request.URL.Hostname())
}

func getSession(c *gin.Context) (db.Session, error) {
	_session := sessions.Default(c)

	session := db.Session{}

	sessionHash, ok := _session.Get("SessionSecretKey").(string)
	if ok == true {
		session, ok = sessionMap[sessionHash]
		if ok {
			return session, nil
		}
	}
	return session, errors.New("не авторизованы")
}

type M struct {
	ID       int32
	Name     string
	Link     string
	Children []M
}

func index(c *gin.Context) {
	session, _ := getSession(c)

	h := gin.H{
		"Title": "Сайтик",
		"Role":  session.User.Role,
	}

	if session.User.Role == "admin" {
		h["JS"] = adminJS
	}

	c.HTML(200, "index", h)
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

	session, _ := getSession(c)
	if session.User.Role != "admin" {
		c.JSON(400, nil)
		return
	}

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
