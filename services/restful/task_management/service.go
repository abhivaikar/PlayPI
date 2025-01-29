package task_management

import (
	"errors"
	"sort"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Due         bool      `json:"due"`
}

var tasks []Task
var taskIDCounter int

func validateTask(task Task) error {
	if len(task.Title) < 3 || len(task.Title) > 100 {
		return errors.New("title must be between 3 and 100 characters")
	}
	if len(task.Description) > 500 {
		return errors.New("description cannot exceed 500 characters")
	}
	if task.Priority != "low" && task.Priority != "medium" && task.Priority != "high" {
		return errors.New("priority must be one of: low, medium, high")
	}
	dueDate, err := time.Parse("2006-01-02", task.DueDate)
	if err != nil {
		return errors.New("due date must follow the format YYYY-MM-DD")
	}
	if dueDate.Before(time.Now()) {
		return errors.New("due date cannot be in the past")
	}
	return nil
}

func CreateTask(newTask Task) (Task, error) {
	if err := validateTask(newTask); err != nil {
		return Task{}, err
	}
	taskIDCounter++
	newTask.ID = taskIDCounter
	newTask.Status = "pending"
	newTask.CreatedAt = time.Now()
	tasks = append(tasks, newTask)
	return newTask, nil
}

func GetTasks() ([]Task, error) {
	if len(tasks) == 0 {
		return nil, errors.New("no tasks created")
	}

	for i := range tasks {
		tasks[i].Due = isTaskDue(tasks[i].DueDate)
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DueDate < tasks[j].DueDate
	})
	return tasks, nil
}

func GetTaskByID(id int) (Task, error) {
	if len(tasks) == 0 {
		return Task{}, errors.New("no tasks created")
	}

	for _, task := range tasks {
		if task.ID == id {
			task.Due = isTaskDue(task.DueDate)
			return task, nil
		}
	}
	return Task{}, errors.New("task not found")
}

func UpdateTask(id int, updatedTask Task) (Task, error) {
	for i, task := range tasks {
		if task.ID == id {
			if err := validateTask(updatedTask); err != nil {
				return Task{}, err
			}
			tasks[i].Title = updatedTask.Title
			tasks[i].Description = updatedTask.Description
			tasks[i].DueDate = updatedTask.DueDate
			tasks[i].Priority = updatedTask.Priority
			tasks[i].Status = updatedTask.Status
			return tasks[i], nil
		}
	}
	return Task{}, errors.New("task not found")
}

func DeleteTask(id int) error {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

func MarkTaskAsCompleted(id int) (Task, error) {
	for i, task := range tasks {
		if task.ID == id {
			if task.Status == "completed" {
				return Task{}, errors.New("task is already marked as completed")
			}
			tasks[i].Status = "completed"
			return tasks[i], nil
		}
	}
	return Task{}, errors.New("task not found")
}

func isTaskDue(dueDate string) bool {
	taskDueDate, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return false
	}
	return taskDueDate.Before(time.Now())
}
