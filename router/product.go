package router

import (
	"awesomeProject2/database"
	"awesomeProject2/utils"
	"github.com/gin-gonic/gin"
)

func getProducts(c *gin.Context) {
	type input struct {
		Search string `json:"Search"`
	}

	i := input{}
	e := c.BindJSON(&i)
	if e != nil {
		utils.Logger.Println(e)
		c.JSON(400, nil)
		return
	}

	c.JSON(200, database.SearchProduct(i.Search))
}

func delProduct(c *gin.Context) {
	session := getSession(c)
	if session.User.Role != "admin" {
		c.JSON(403, nil)
		return
	}

	type input struct {
		ID int `json:"id"`
	}

	i := input{}
	e := c.BindJSON(&i)
	if e != nil {
		utils.Logger.Println(e)
		c.JSON(400, nil)
		return
	}

	e = database.DeleteProduct(i.ID)
	if e != nil {
		c.JSON(400, nil)
		return
	}

	c.JSON(200, nil)
}

func updateCart(c *gin.Context) {
	session := getSession(c)
	if !session.Exists {
		c.JSON(403, nil)
		return
	}

	type input struct {
		ID    int
		Count int
	}

	i := input{}
	e := c.BindJSON(&i)
	if e != nil {
		utils.Logger.Println(e)
		c.JSON(400, nil)
		return
	}

	e = database.UpdateCart(session.User.Login, i.ID, i.Count)
	if e != nil {
		c.JSON(400, nil)
		return
	}

	count, e := database.SelectCart(session.User.Login)
	if e != nil {
		c.JSON(400, nil)
		return
	}

	c.JSON(200, count)
}

func selectCart(c *gin.Context) {
	session := getSession(c)
	if !session.Exists {
		c.JSON(403, nil)
		return
	}

	data, e := database.SelectCart(session.User.Login)
	if e != nil {
		c.JSON(400, nil)
		return
	}

	c.JSON(200, data)
}
