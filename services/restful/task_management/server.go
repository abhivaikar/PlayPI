package task_management

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := setupRouter()
	r.Run(":8085")
}

func StartServerForTesting() *gin.Engine {
	return setupRouter()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/tasks", func(c *gin.Context) {
		var newTask Task
		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task, err := CreateTask(newTask)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, task)
	})

	r.GET("/tasks", func(c *gin.Context) {
		tasks, err := GetTasks()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tasks)
	})

	r.GET("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
			return
		}
		task, err := GetTaskByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, task)
	})

	r.PUT("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
			return
		}
		var updatedTask Task
		if err := c.ShouldBindJSON(&updatedTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task, err := UpdateTask(id, updatedTask)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, task)
	})

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
			return
		}
		if err := DeleteTask(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
	})

	r.PUT("/tasks/:id/complete", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
			return
		}
		task, err := MarkTaskAsCompleted(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, task)
	})

	return r
}
