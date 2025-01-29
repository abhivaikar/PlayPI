package task_management

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupTestServer() *gin.Engine {
	// Start the server without any pre-seeded tasks
	tasks = []Task{}
	taskIDCounter = 0
	return StartServerForTesting()
}

func TestCreateTask(t *testing.T) {
	r := setupTestServer()

	t.Run("Create a valid task", func(t *testing.T) {
		futureDate := getFutureDate(30) // Set due date 30 days from now

		payload := `{
			"title": "New Task",
			"description": "Task description",
			"due_date": "` + futureDate + `",
			"priority": "high"
		}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusCreated, resp.Code)

		var task Task
		err := json.Unmarshal(resp.Body.Bytes(), &task)
		require.NoError(t, err)
		require.Equal(t, "New Task", task.Title)
		require.Equal(t, "Task description", task.Description)
		require.Equal(t, futureDate, task.DueDate)
		require.Equal(t, "high", task.Priority)
		require.Equal(t, "pending", task.Status)
	})

	// Validation cases
	t.Run("Validation Error - Missing Title", func(t *testing.T) {
		payload := `{
			"description": "Task description",
			"due_date": "2025-02-01",
			"priority": "medium"
		}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "title must be between 3 and 100 characters", response["error"])
	})

	t.Run("Validation Error - Title Too Short", func(t *testing.T) {
		payload := `{
			"title": "Up",
			"description": "Task description",
			"due_date": "2025-02-01",
			"priority": "medium"
		}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "title must be between 3 and 100 characters", response["error"])
	})

	t.Run("Validation Error - Description Too Long", func(t *testing.T) {
		payload := `{
			"title": "Task with Long Description",
			"description": "` + generateLongString(501) + `",
			"due_date": "2025-02-01",
			"priority": "medium"
		}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "description cannot exceed 500 characters", response["error"])
	})

	t.Run("Validation Error - Invalid Due Date Format", func(t *testing.T) {
		payload := `{
			"title": "Task with Invalid Due Date",
			"description": "Task description",
			"due_date": "31-01-2025",
			"priority": "medium"
		}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "due date must follow the format YYYY-MM-DD", response["error"])
	})

	t.Run("Validation Error - Past Due Date", func(t *testing.T) {
		payload := `{
			"title": "Task with Past Due Date",
			"description": "Task description",
			"due_date": "2022-01-01",
			"priority": "medium"
		}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "due date cannot be in the past", response["error"])
	})

	t.Run("Validation Error - Invalid Priority", func(t *testing.T) {
		payload := `{
			"title": "Task with Invalid Priority",
			"description": "Task description",
			"due_date": "2025-02-01",
			"priority": "urgent"
		}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "priority must be one of: low, medium, high", response["error"])
	})
}

func TestGetTasks(t *testing.T) {
	r := setupTestServer()

	t.Run("No tasks available", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "no tasks created", response["message"])
	})

	t.Run("Get all tasks", func(t *testing.T) {
		// Create a new task to ensure tasks are available
		payload := `{
			"title": "Test Task",
			"description": "This is a test task",
			"due_date": "2025-01-31",
			"priority": "medium"
		}`
		reqCreate, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		reqCreate.Header.Set("Content-Type", "application/json")
		respCreate := httptest.NewRecorder()
		r.ServeHTTP(respCreate, reqCreate)

		req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		var tasks []Task
		err := json.Unmarshal(resp.Body.Bytes(), &tasks)
		require.NoError(t, err)
		require.Len(t, tasks, 1)
		require.Equal(t, "Test Task", tasks[0].Title)
		require.Equal(t, "This is a test task", tasks[0].Description)
	})
}

func TestGetTaskByID(t *testing.T) {
	r := setupTestServer()

	t.Run("No tasks available", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "no tasks created", response["error"])
	})

	t.Run("Get task by valid ID", func(t *testing.T) {
		// Create a new task to ensure tasks are available
		payload := `{
			"title": "Test Task",
			"description": "This is a test task",
			"due_date": "2025-01-31",
			"priority": "medium"
		}`
		reqCreate, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		reqCreate.Header.Set("Content-Type", "application/json")
		respCreate := httptest.NewRecorder()
		r.ServeHTTP(respCreate, reqCreate)

		req, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		var task Task
		err := json.Unmarshal(resp.Body.Bytes(), &task)
		require.NoError(t, err)
		require.Equal(t, "Test Task", task.Title)
		require.Equal(t, "This is a test task", task.Description)
		require.Equal(t, "2025-01-31", task.DueDate)
	})

	t.Run("Get task by invalid ID", func(t *testing.T) {
		// Create a new task to ensure tasks are available
		payload := `{
			"title": "Test Task",
			"description": "This is a test task",
			"due_date": "2025-01-31",
			"priority": "medium"
		}`
		reqCreate, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		reqCreate.Header.Set("Content-Type", "application/json")
		respCreate := httptest.NewRecorder()
		r.ServeHTTP(respCreate, reqCreate)

		req, _ := http.NewRequest(http.MethodGet, "/tasks/999", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "task not found", response["error"])
	})
}

func TestUpdateTask(t *testing.T) {
	r := setupTestServer()

	t.Run("Update an existing task with valid data", func(t *testing.T) {
		taskID := createTaskForTest(t, r)   // Create task and get its ID
		updatedDueDate := getFutureDate(60) // 60 days from now

		payload := `{
			"title": "Updated Task",
			"description": "Updated description",
			"due_date": "` + updatedDueDate + `",
			"priority": "high"
		}`
		reqUpdate, _ := http.NewRequest(http.MethodPut, "/tasks/"+strconv.Itoa(taskID), bytes.NewBufferString(payload))
		reqUpdate.Header.Set("Content-Type", "application/json")
		respUpdate := httptest.NewRecorder()
		r.ServeHTTP(respUpdate, reqUpdate)

		require.Equal(t, http.StatusOK, respUpdate.Code)

		var task Task
		err := json.Unmarshal(respUpdate.Body.Bytes(), &task)
		require.NoError(t, err)
		require.Equal(t, "Updated Task", task.Title)
		require.Equal(t, "Updated description", task.Description)
		require.Equal(t, updatedDueDate, task.DueDate)
		require.Equal(t, "high", task.Priority)
	})

	// Validation Tests
	t.Run("Validation Error - Missing Title", func(t *testing.T) {
		taskID := createTaskForTest(t, r) // Create task and get its ID
		payload := `{
			"title": "",
			"description": "Updated description",
			"due_date": "` + getFutureDate(30) + `",
			"priority": "medium"
		}`
		reqUpdate, _ := http.NewRequest(http.MethodPut, "/tasks/"+strconv.Itoa(taskID), bytes.NewBufferString(payload))
		reqUpdate.Header.Set("Content-Type", "application/json")
		respUpdate := httptest.NewRecorder()
		r.ServeHTTP(respUpdate, reqUpdate)

		require.Equal(t, http.StatusBadRequest, respUpdate.Code)

		var response map[string]string
		err := json.Unmarshal(respUpdate.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "title must be between 3 and 100 characters", response["error"])
	})

	t.Run("Validation Error - Invalid Due Date Format", func(t *testing.T) {
		taskID := createTaskForTest(t, r) // Create task and get its ID
		payload := `{
			"title": "Valid Title",
			"description": "Valid description",
			"due_date": "31-01-2025",
			"priority": "medium"
		}`
		reqUpdate, _ := http.NewRequest(http.MethodPut, "/tasks/"+strconv.Itoa(taskID), bytes.NewBufferString(payload))
		reqUpdate.Header.Set("Content-Type", "application/json")
		respUpdate := httptest.NewRecorder()
		r.ServeHTTP(respUpdate, reqUpdate)

		require.Equal(t, http.StatusBadRequest, respUpdate.Code)

		var response map[string]string
		err := json.Unmarshal(respUpdate.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "due date must follow the format YYYY-MM-DD", response["error"])
	})

	t.Run("Validation Error - Past Due Date", func(t *testing.T) {
		taskID := createTaskForTest(t, r) // Create task and get its ID
		payload := `{
			"title": "Valid Title",
			"description": "Valid description",
			"due_date": "2020-01-01",
			"priority": "medium"
		}`
		reqUpdate, _ := http.NewRequest(http.MethodPut, "/tasks/"+strconv.Itoa(taskID), bytes.NewBufferString(payload))
		reqUpdate.Header.Set("Content-Type", "application/json")
		respUpdate := httptest.NewRecorder()
		r.ServeHTTP(respUpdate, reqUpdate)

		require.Equal(t, http.StatusBadRequest, respUpdate.Code)

		var response map[string]string
		err := json.Unmarshal(respUpdate.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "due date cannot be in the past", response["error"])
	})

	t.Run("Validation Error - Invalid Priority", func(t *testing.T) {
		taskID := createTaskForTest(t, r) // Create task and get its ID
		payload := `{
			"title": "Valid Title",
			"description": "Valid description",
			"due_date": "` + getFutureDate(30) + `",
			"priority": "urgent"
		}`
		reqUpdate, _ := http.NewRequest(http.MethodPut, "/tasks/"+strconv.Itoa(taskID), bytes.NewBufferString(payload))
		reqUpdate.Header.Set("Content-Type", "application/json")
		respUpdate := httptest.NewRecorder()
		r.ServeHTTP(respUpdate, reqUpdate)

		require.Equal(t, http.StatusBadRequest, respUpdate.Code)

		var response map[string]string
		err := json.Unmarshal(respUpdate.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "priority must be one of: low, medium, high", response["error"])
	})

	t.Run("Validation Error - Description Too Long", func(t *testing.T) {
		taskID := createTaskForTest(t, r) // Create task and get its ID
		payload := `{
			"title": "Valid Title",
			"description": "` + generateLongString(501) + `",
			"due_date": "` + getFutureDate(30) + `",
			"priority": "medium"
		}`
		reqUpdate, _ := http.NewRequest(http.MethodPut, "/tasks/"+strconv.Itoa(taskID), bytes.NewBufferString(payload))
		reqUpdate.Header.Set("Content-Type", "application/json")
		respUpdate := httptest.NewRecorder()
		r.ServeHTTP(respUpdate, reqUpdate)

		require.Equal(t, http.StatusBadRequest, respUpdate.Code)

		var response map[string]string
		err := json.Unmarshal(respUpdate.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "description cannot exceed 500 characters", response["error"])
	})
}

func TestDeleteTask(t *testing.T) {
	r := setupTestServer()

	t.Run("Delete an existing task", func(t *testing.T) {
		// First, create a task
		payload := `{
			"title": "Task to Delete",
			"description": "Description",
			"due_date": "2025-02-01",
			"priority": "medium"
		}`
		reqCreate, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		reqCreate.Header.Set("Content-Type", "application/json")
		respCreate := httptest.NewRecorder()
		r.ServeHTTP(respCreate, reqCreate)

		// Delete the task
		reqDelete, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
		respDelete := httptest.NewRecorder()
		r.ServeHTTP(respDelete, reqDelete)

		require.Equal(t, http.StatusOK, respDelete.Code)

		var response map[string]string
		err := json.Unmarshal(respDelete.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "task deleted", response["message"])
	})

	t.Run("Delete non-existing task", func(t *testing.T) {
		reqDelete, _ := http.NewRequest(http.MethodDelete, "/tasks/999", nil)
		respDelete := httptest.NewRecorder()
		r.ServeHTTP(respDelete, reqDelete)

		require.Equal(t, http.StatusNotFound, respDelete.Code)

		var response map[string]string
		err := json.Unmarshal(respDelete.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "task not found", response["error"])
	})
}

func TestMarkTaskAsCompleted(t *testing.T) {
	r := setupTestServer()

	t.Run("Mark task as completed", func(t *testing.T) {
		// First, create a task
		payload := `{
			"title": "Task to Complete",
			"description": "Description",
			"due_date": "2025-02-01",
			"priority": "medium"
		}`
		reqCreate, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		reqCreate.Header.Set("Content-Type", "application/json")
		respCreate := httptest.NewRecorder()
		r.ServeHTTP(respCreate, reqCreate)

		// Mark task as completed
		reqComplete, _ := http.NewRequest(http.MethodPut, "/tasks/1/complete", nil)
		respComplete := httptest.NewRecorder()
		r.ServeHTTP(respComplete, reqComplete)

		require.Equal(t, http.StatusOK, respComplete.Code)

		var task Task
		err := json.Unmarshal(respComplete.Body.Bytes(), &task)
		require.NoError(t, err)
		require.Equal(t, "completed", task.Status)
	})

	t.Run("Mark non-existing task as completed", func(t *testing.T) {
		reqComplete, _ := http.NewRequest(http.MethodPut, "/tasks/999/complete", nil)
		respComplete := httptest.NewRecorder()
		r.ServeHTTP(respComplete, reqComplete)

		require.Equal(t, http.StatusNotFound, respComplete.Code)

		var response map[string]string
		err := json.Unmarshal(respComplete.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "task not found", response["error"])
	})
}

func generateLongString(length int) string {
	return string(bytes.Repeat([]byte("a"), length))
}

func getFutureDate(daysFromNow int) string {
	// Calculate the date after `daysFromNow` days
	futureDate := time.Now().AddDate(0, 0, daysFromNow)
	// Format the date as YYYY-MM-DD
	return futureDate.Format("2006-01-02")
}

func createTaskForTest(t *testing.T, r *gin.Engine) int {
	futureDate := getFutureDate(30) // 30 days from now
	payload := `{
		"title": "Task for Update Test",
		"description": "Description for update test",
		"due_date": "` + futureDate + `",
		"priority": "medium"
	}`

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)

	var task Task
	err := json.Unmarshal(resp.Body.Bytes(), &task)
	require.NoError(t, err)

	return task.ID
}
