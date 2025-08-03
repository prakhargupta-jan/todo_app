package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const filename string = "todos.txt"

func GetTodos() ([]Todo, error) {
	var data Todos
	if err := ReadJSON(filename, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func checkEr(c *gin.Context, err error, status int) bool {
	if err != nil {
		log.Println("ERROR : ", err.Error())
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return true
	}
	return false
}

func main() {
	router := gin.Default()

	// Get All Todos
	router.GET("/", func(c *gin.Context) {
		todos, err := GetTodos()
		if checkEr(c, err, http.StatusBadGateway) {
			return
		}
		c.JSON(http.StatusOK, todos)
	})
	// Post TODO
	router.POST("/", func(c *gin.Context) {
		todo := Todo{}
		err := c.ShouldBindBodyWithJSON(&todo)
		if checkEr(c, err, http.StatusBadRequest) {
			return
		}

		todos, err := GetTodos()
		if checkEr(c, err, http.StatusBadGateway) {
			return
		}

		todos = append(todos, todo)
		err = WriteJSON(filename, todos)
		if checkEr(c, err, http.StatusBadGateway) {
			return
		}
		c.JSON(http.StatusOK, todo)
	})

	// Update TODO
	router.PUT("/:idx", func(c *gin.Context) {
		newTodo := UpdateTodo{}

		err := c.ShouldBindBodyWithJSON(&newTodo)
		if checkEr(c, err, http.StatusBadRequest) {
			return
		}

		todos, err := GetTodos()
		if checkEr(c, err, http.StatusBadGateway) {
			return
		}

		uriParams := struct {
			Idx int `uri:"idx"`
		}{}
		err = c.ShouldBindUri(&uriParams)
		if checkEr(c, err, http.StatusBadRequest) {
			return
		}
		if uriParams.Idx < 0 || uriParams.Idx >= len(todos) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("No todo found with id %d.", uriParams.Idx),
			})
			return
		}

		UpdateFields(&todos[uriParams.Idx], &newTodo)
		err = WriteJSON(filename, todos)
		if checkEr(c, err, http.StatusBadGateway) {
			return
		}
		c.JSON(http.StatusOK, todos[uriParams.Idx])
	})

	// Get TODO
	router.GET("/:idx", func(c *gin.Context) {
		uriParams := struct {
			Idx int `uri:"idx"`
		}{}
		err := c.ShouldBindUri(&uriParams)
		if checkEr(c, err, http.StatusBadRequest) {
			return
		}
		todos, err := GetTodos()
		if checkEr(c, err, http.StatusBadGateway) {
			return
		}
		if uriParams.Idx < 0 || uriParams.Idx >= len(todos) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("No todo found with id %d.", uriParams.Idx),
			})
			return
		}
		todo := todos[uriParams.Idx]
		c.JSON(http.StatusOK, todo)
	})

	// Delete TODO
	router.DELETE("/:idx", func(c *gin.Context) {
		todos, err := GetTodos()
		if checkEr(c, err, http.StatusBadGateway) {
			return
		}

		uriParams := struct {
			Idx int `uri:"idx"`
		}{}
		err = c.ShouldBindUri(&uriParams)
		if checkEr(c, err, http.StatusBadRequest) {
			return
		}
		if uriParams.Idx < 0 || uriParams.Idx >= len(todos) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("No todo found with id %d.", uriParams.Idx),
			})
			return
		}
		log.Println(todos)
		todos = append(todos[:uriParams.Idx], todos[uriParams.Idx+1:]...)
		log.Println(todos)
		if err := WriteJSON(filename, todos); err != nil {
			log.Println("ERROR: C", err)
			c.JSON(http.StatusBadGateway, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Todo deleted successfully.",
		})
	})

	// Toggle TODO
	router.PUT("/:idx/toggle", func(c *gin.Context) {
		todos, err := GetTodos()
		if checkEr(c, err, http.StatusBadGateway) {
			return
		}

		uriParams := struct {
			Idx int `uri:"idx"`
		}{}
		err = c.ShouldBindUri(&uriParams)
		if checkEr(c, err, http.StatusBadRequest) {
			return
		}

		if uriParams.Idx < 0 || uriParams.Idx >= len(todos) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("No todo found with id %d.", uriParams.Idx),
			})
			return
		}
		todos[uriParams.Idx].Toggle()
		err = WriteJSON(filename, todos)
		if checkEr(c, err, http.StatusBadGateway) {
			return
		}
		c.JSON(http.StatusOK, todos[uriParams.Idx])
	})

	router.Run()
}
