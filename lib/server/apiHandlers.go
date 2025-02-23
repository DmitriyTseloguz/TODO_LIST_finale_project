package server

import (
	"net/http"
	"time"

	planner "todo-list/lib"
	"todo-list/lib/extensions"
)

var ApiHandlers = map[string]http.HandlerFunc{
	"/api/nextdate": func(response http.ResponseWriter, request *http.Request) {
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

		var validationError = repeater.Validate()

		if validationError != nil {
			response.Write([]byte(validationError.Error()))
		}

		var event = planner.NewEvent(date, repeater)

		var rescheduler, err = planner.DefineRescheduler(event)

		if err != nil {
			response.Write([]byte(err.Error()))

			return
		}

		if nowParseError == nil {
			rescheduler.SetBaseOnDate(baseOnDate)
		}

		rescheduler.Reschedule(event)

		var nextTime = event.GetTime()

		response.Write([]byte(nextTime.Format("20060102")))
	},
}
