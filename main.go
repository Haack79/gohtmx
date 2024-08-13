package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	name := "World!"
	fmt.Printf("hello, %s!\n", name)
	// initialize db
	InitDatabase()
	defer DB.Close()

	e := gin.Default()
	e.LoadHTMLGlob("templates/*")

	e.GET("/", func(c *gin.Context) {
		todos := ReadToDoList()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"todos": todos,
			"name":  name,
		})
	})

	e.POST("/todos", func(c *gin.Context) {
		title := c.PostForm("title")
		status := c.PostForm("status")
		_, err := CreateToDo(title, status)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to create todo: %v", err)
			return
		}

		// Re-fetch the updated todos
		todos := ReadToDoList()

		// Render only the tasks partial
		c.HTML(http.StatusOK, "tasks_partial.html", gin.H{
			"todos": todos,
		})
	})

	e.DELETE("/todos/:id", func(c *gin.Context) {
		param := c.Param("id")
		id, _ := strconv.ParseInt(param, 10, 64)
		DeleteToDo(id)
	})

	e.Run(":8080")
}
