package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"todo-list/lib/database"
	"todo-list/lib/extensions"
	"todo-list/lib/model"
	"todo-list/lib/planner"
	"todo-list/lib/utils"
)

var ApiHandlers = map[string]http.HandlerFunc{
	"/api/nextdate": nextDate,
	"/api/task": func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getTask(response, request)
		case http.MethodPost:
			createTask(response, request)
		case http.MethodPut:
			editTask(response, request)
		default:
			http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		}
	},
	"/api/tasks": func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getTasks(response, request)
		default:
			http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		}
	},
}

func createTask(response http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Request body is empty"})
		return
	}

	body, err := io.ReadAll(request.Body)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Failed to read request body"})
		return
	}

	defer request.Body.Close()

	if len(body) == 0 {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Request body is empty"})
		return
	}

	if !json.Valid(body) {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	var task, taskDecodeError = utils.DecodeTaskFromJSON(body)

	if taskDecodeError != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": taskDecodeError.Error()})
		return
	}

	db := database.GetDB()
	id, err := db.CreateTask(task)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(map[string]string{"error": "Failed to create task"})
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]int{"id": id})
}

func getTask(response http.ResponseWriter, request *http.Request) {
	var id = request.URL.Query().Get("id")

	if id == "" {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Missing 'id' parameter"})
		return
	}

	taskID, err := strconv.Atoi(id)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(map[string]string{"error": "Invalid 'id' parameter"})
		return
	}

	db := database.GetDB()
	task, err := db.GetTask(taskID)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(map[string]string{"error": fmt.Sprintf("Failed to get task: %v", err)})
		return
	}

	// Возвращаем задачу в формате JSON
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(task)
}

func getTasks(response http.ResponseWriter, request *http.Request) {
	var db = database.GetDB()

	var tasks, err = db.GetAllTasks()

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(map[string]string{"error": "Failed to get tasks"})
		return
	}

	var responseTasks = struct {
		Tasks []model.Task `json:"tasks"`
	}{
		Tasks: tasks,
	}

	var jsonTasks, encodeError = json.Marshal(responseTasks)

	fmt.Println("JSON: " + string(jsonTasks))
	if encodeError != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": encodeError.Error()})
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonTasks)
}

func editTask(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	if request.Body == nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Request body is empty"})
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Failed to read request body"})
		return
	}
	defer request.Body.Close()

	if len(body) == 0 {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Request body is empty"})
		return
	}

	if !json.Valid(body) {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	task, err := utils.DecodeTaskFromJSON(body)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": err.Error()})
		return
	}

	if task.ID == "" {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Missing 'id' field"})
		return
	}

	taskID, err := strconv.Atoi(task.ID)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{"error": "Invalid 'id' field"})
		return
	}

	db := database.GetDB()
	err = db.UpdateTask(taskID, task)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(map[string]string{"error": fmt.Sprintf("Failed to update task: %v", err)})
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(map[string]interface{}{})
}

func nextDate(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var nowParameter = request.URL.Query().Get("now")
	var dateParameter = request.URL.Query().Get("date")
	var repeatParameter = request.URL.Query().Get("repeat")

	var baseOnDate, nowParseError = time.Parse("20060102", nowParameter)
	var date, dateParseError = time.Parse("20060102", dateParameter)
	var repeater = extensions.RepeaterRule(repeatParameter)

	if dateParseError != nil {
		response.Write([]byte("Wrong date"))
		return
	}

	if dateParameter == "" {
		response.Write([]byte(""))
		return
	}

	if nowParseError != nil {
		baseOnDate = time.Now()
	}

	var validationError = repeater.Validate()

	if validationError != nil {
		response.Write([]byte(validationError.Error()))
		return
	}

	var event = planner.NewEvent(date, repeater)
	var nextDate, err = event.GetNextDate(baseOnDate)

	if err != nil {
		response.Write([]byte(err.Error()))

		return
	}

	response.Write([]byte(nextDate.Format("20060102")))
}
