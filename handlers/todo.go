package handlers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gocql/gocql"
    "todo-api/db"
    "todo-api/models"
)

func CreateTodoHandler(c *gin.Context) {
    var todo models.Todo
    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    todo.ID = gocql.TimeUUID()
    todo.Created = time.Now()
    todo.Updated = time.Now()

    if err := db.Session.Query(`INSERT INTO todos (id, user_id, title, description, status, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)`,
        todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated).Exec(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, todo)
}

func ListAllTodosHandler(c *gin.Context) {
    var todos []models.Todo

    sortBy := c.DefaultQuery("sortBy", "created") 

    query := `SELECT id, user_id, title, description, status, created, updated FROM todos`
    query += ` ORDER BY ` + sortBy + ` DESC`
    iter := db.Session.Query(query).Iter()

    var todo models.Todo
    for iter.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
        todos = append(todos, todo)
    }

    if err := iter.Close(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, todos)
}


func GetTodoHandler(c *gin.Context) {
    userID, err := gocql.ParseUUID(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    todoID, err := gocql.ParseUUID(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
        return
    }

    var todo models.Todo

    if err := db.Session.Query(`SELECT id, user_id, title, description, status, created, updated FROM todos WHERE user_id = ? AND id = ?`,
        userID, todoID).Consistency(gocql.One).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }

    c.JSON(http.StatusOK, todo)
}

func UpdateTodoHandler(c *gin.Context) {
    userID, err := gocql.ParseUUID(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    todoID, err := gocql.ParseUUID(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
        return
    }

    var todo models.Todo
    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    todo.Updated = time.Now()

    if err := db.Session.Query(`UPDATE todos SET title = ?, description = ?, status = ?, updated = ? WHERE user_id = ? AND id = ?`,
        todo.Title, todo.Description, todo.Status, todo.Updated, userID, todoID).Exec(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, todo)
}

func DeleteTodoHandler(c *gin.Context) {
    userID, err := gocql.ParseUUID(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    todoID, err := gocql.ParseUUID(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
        return
    }

    if err := db.Session.Query(`DELETE FROM todos WHERE user_id = ? AND id = ?`, userID, todoID).Exec(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

func ListTodosHandler(c *gin.Context) {
    userID, err := gocql.ParseUUID(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    status := c.Query("status")
    pageSize := 10 
    page := c.Query("page")
    sortBy := c.DefaultQuery("sortBy", "created") 

    pageNumber, err := strconv.Atoi(page)
    if err != nil || pageNumber < 1 {
        pageNumber = 1
    }

    offset := (pageNumber - 1) * pageSize

    var todos []models.Todo
    var iter *gocql.Iter

    query := `SELECT id, user_id, title, description, status, created, updated FROM todos WHERE user_id = ?`
    if status != "" {
        query += ` AND status = ?`
        query += ` ORDER BY ` + sortBy + ` DESC LIMIT ? OFFSET ?`
        iter = db.Session.Query(query, userID, status, pageSize, offset).Iter()
    } else {
        query += ` ORDER BY ` + sortBy + ` DESC LIMIT ? OFFSET ?`
        iter = db.Session.Query(query, userID, pageSize, offset).Iter()
    }

    var todo models.Todo
    for iter.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
        todos = append(todos, todo)
    }

    if err := iter.Close(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, todos)
}
