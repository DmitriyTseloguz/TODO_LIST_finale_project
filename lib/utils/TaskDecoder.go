package utils

import (
	"encoding/json"
	"fmt"
	"time"
	"todo-list/lib/model"
	"todo-list/lib/planner"
)

func DecodeTaskFromJSON(body []byte) (*model.Task, error) {
	var task model.Task

	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if task.Title == "" {
		return nil, fmt.Errorf("не указан заголовок задачи")
	}

	var repeaterError = task.GetRepeater().Validate()

	if len(task.GetRepeater()) > 0 && repeaterError != nil {
		return nil, repeaterError
	}

	var now = time.Now().Truncate(24 * time.Hour)

	if task.GetTime().Before(now) {
		task.SetTime(now)

		var event = planner.NewEvent(task.GetTime(), task.GetRepeater())
		var nextDate, err = event.GetNextDate(now)

		if err != nil && task.GetRepeater() != "" {
			return nil, err
		}

		if task.GetRepeater() != "" {
			task.SetTime(nextDate)
		}
	}

	return &task, nil
}
