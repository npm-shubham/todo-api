package main

import (
    "github.com/gin-gonic/gin"
    "todo-api/db"
    "todo-api/handlers"
)

func main() {
    db.InitScyllaDB()
    defer db.Session.Close()

    router := gin.Default()

    router.POST("/todos", handlers.CreateTodoHandler)
    router.GET("/todos/:user_id/:id", handlers.GetTodoHandler)
    router.PUT("/todos/:user_id/:id", handlers.UpdateTodoHandler)
    router.DELETE("/todos/:user_id/:id", handlers.DeleteTodoHandler)
    router.GET("/todos/:user_id", handlers.ListTodosHandler)

    router.Run(":8080")
}
