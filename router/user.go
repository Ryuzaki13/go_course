package router

import (
	"awesomeProject2/database"
	"awesomeProject2/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func reg(c *gin.Context) {
	var user database.User
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

	var user database.User
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
		hash, ok := database.CreateSession(&user)
		if ok {
			session.Set("SessionSecretKey", hash)
			e = session.Save()
			if e != nil {
				utils.Logger.Println(e)
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
