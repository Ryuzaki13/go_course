package router

import (
	"awesomeProject2/database"
	"awesomeProject2/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
)

var adminJS template.JS
var adminCSS template.CSS

func Initialized() *gin.Engine {
	// Загрузить в память админский функционал
	adminJS = template.JS(utils.LoadAssets("secret/js/admin.js"))
	adminCSS = template.CSS(utils.LoadAssets("secret/css/admin.css"))

	// Создать куки
	word := sessions.NewCookieStore([]byte("SecretX"))

	router := gin.Default()

	// Записать куки в сессию
	router.Use(sessions.Sessions("session", word))

	router.LoadHTMLGlob("template/*")
	router.Static("assets", "assets")

	router.GET("/", index)
	router.GET("/cart", handlerCart)

	routerUser := router.Group("/user")

	routerUser.POST("/reg", reg)
	routerUser.POST("/login", login)
	routerUser.POST("/logout", logout)

	routerAPI := router.Group("/api")

	routerAPI.POST("/product", getProducts)
	routerAPI.DELETE("/product", delProduct)

	routerAPI.POST("/cart/update", updateCart)
	routerAPI.GET("/cart", selectCart)

	return router
}

func getSession(c *gin.Context) *database.Session {
	_session := sessions.Default(c)

	sessionHash, ok := _session.Get("SessionSecretKey").(string)
	if ok {
		session := database.GetSession(sessionHash)
		if session != nil {
			session.Exists = true
			return session
		}
	}

	return &database.Session{
		Exists: false,
	}
}

func index(c *gin.Context) {
	session := getSession(c)

	h := gin.H{
		"Title": "Сайтик",
		"Role":  session.User.Role,
	}

	if session.User.Role == "admin" {
		h["JS"] = adminJS
	}

	c.HTML(200, "index", h)
}

func handlerCart(c *gin.Context) {
	session := getSession(c)

	cart, _ := database.SelectCart(session.User.Login)

	h := gin.H{
		"Title": "Сайтик",
		"Cart":  cart,
		"Role":  session.User.Role,
	}

	if session.User.Role == "admin" {
		h["JS"] = adminJS
	}

	c.HTML(200, "cart", h)
}
