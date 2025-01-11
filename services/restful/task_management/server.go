package task_management

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"` // Format: YYYY-MM-DD
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Due         bool      `json:"due"` // New field to indicate if the task is overdue
}

var tasks []Task
var taskIDCounter int

func StartServer() {
	r := gin.Default()

	// Create a Task
	r.POST("/tasks", func(c *gin.Context) {
		var newTask Task
		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		taskIDCounter++
		newTask.ID = taskIDCounter
		newTask.Status = "pending"
		newTask.CreatedAt = time.Now()
		tasks = append(tasks, newTask)
		c.JSON(http.StatusCreated, newTask)
	})

	// Get All Tasks
	r.GET("/tasks", func(c *gin.Context) {
		updatedTasks := make([]Task, len(tasks))
		for i, task := range tasks {
			task.Due = isTaskDue(task.DueDate)
			updatedTasks[i] = task
		}
		c.JSON(http.StatusOK, updatedTasks)
	})

	// Get a Task by ID
	r.GET("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}
		for _, task := range tasks {
			if task.ID == id {
				task.Due = isTaskDue(task.DueDate)
				c.JSON(http.StatusOK, task)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	})

	// Update a Task
	r.PUT("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}
		var updatedTask Task
		if err := c.ShouldBindJSON(&updatedTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Title = updatedTask.Title
				tasks[i].Description = updatedTask.Description
				tasks[i].DueDate = updatedTask.DueDate
				tasks[i].Priority = updatedTask.Priority
				tasks[i].Status = updatedTask.Status
				c.JSON(http.StatusOK, tasks[i])
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	})

	// Delete a Task
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}
		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	})

	// Mark a Task as Completed
	r.PUT("/tasks/:id/complete", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Status = "completed"
				c.JSON(http.StatusOK, tasks[i])
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	})

	r.Run(":8085") // Task management API runs on port 8085
}

func isTaskDue(dueDate string) bool {
	// Parse the due date
	taskDueDate, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		log.Printf("Error parsing due date: %v", err)
		return false
	}
	// Compare with the current date
	return taskDueDate.Before(time.Now())
}
