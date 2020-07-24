package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	controllers "github.com/jayschoen/iWant-backend/controllers"
	helpers "github.com/jayschoen/iWant-backend/helpers"
)

func get(w http.ResponseWriter, r *http.Request) {

	requestBody := helpers.ParseSlackPayload(r)

	if len(requestBody) > 0 {

		cmdText := helpers.ParseSlackPayloadText(requestBody["text"][0])

		// empty slices will always have len() of 1 apparently, so check for empty string
		if cmdText[0] == "" {

			wants, err := controllers.GetAllWants()
			if err != nil {
				helpers.RespondWithError(w, helpers.ArgParser(http.StatusInternalServerError, err.Error()))
				return
			}

			helpers.RespondWithJSON(w, helpers.ArgParser(http.StatusOK, wants))
			return

		}

		if len(cmdText) >= 1 {

			id, err := strconv.Atoi(cmdText[0])
			if err != nil || id < 1 {
				helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Invalid Want ID"))
				return
			}

			want, err := controllers.GetWantByID(id)
			if err != nil {
				helpers.RespondWithError(w, helpers.ArgParser(http.StatusInternalServerError, err.Error()))
				return
			}

			helpers.RespondWithJSON(w, helpers.ArgParser(http.StatusOK, want))
			return
		}
	}
	helpers.RespondWithError(w, helpers.ArgParser(http.StatusUnprocessableEntity, "Missing Request Body"))
	return
}

// TODO: NEEDS A DIALOG POPUP IN SLACK
// TODO:
func post(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	params := r.URL.Query()

	paramID, hasID := vars["id"]

	if !hasID {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Missing Want ID"))
		return
	}

	id, err := strconv.Atoi(paramID)
	if err != nil || id < 1 {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Invalid Want ID"))
		return
	}

	if len(params) == 0 {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Parameters Are Required"))
		return
	}

	status, statusExists := params["status"]
	wants, wantsExists := params["wants"]
	targetTime, targetTimeExists := params["targetTime"]

	if !statusExists && !wantsExists && !targetTimeExists {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Missing Parameter - Request Must Include At Least status Or wants Or targetTime"))
		return
	}

	if err := controllers.UpdateWant(id, status, wants, targetTime); err != nil {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.RespondWithJSON(w, helpers.ArgParser(http.StatusOK, map[string]interface{}{"result": "success"}))
	return

}

// TODO: NEEDS A DIALOG POPUP IN SLACK
// TODO: AND NEEDS TO GET SLACK USERS ID
// TODO: AND NEEDS TO HANDLE targetTime PROPERLY
func put(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	if len(params) == 0 {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Parameters Are Required"))
		return
	}

	paramSlackID, slackIDExists := params["slackID"]
	status, statusExists := params["status"]
	wants, wantsExists := params["wants"]
	targetTime, targetTimeExists := params["targetTime"]

	exists := make(map[string]bool)
	exists["slackID"] = slackIDExists
	exists["status"] = statusExists
	exists["wants"] = wantsExists
	exists["targetTime"] = targetTimeExists

	for paramName, exist := range exists {
		if !exist {
			helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, fmt.Sprintf("Missing Parameter: %s", paramName)))
			return
		}
	}

	slackID, err := strconv.Atoi(paramSlackID[0])
	if err != nil || slackID < 1 {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Invalid Slack ID"))
		return
	}

	if err := controllers.InsertWant(slackID, status[0], wants[0], helpers.ParseTimeString(targetTime[0])); err != nil {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.RespondWithJSON(w, helpers.ArgParser(http.StatusOK, map[string]interface{}{"result": "success"}))
	return
}

func delete(w http.ResponseWriter, r *http.Request) {

	requestBody := helpers.ParseSlackPayload(r)

	if len(requestBody) > 0 {

		cmdText := helpers.ParseSlackPayloadText(requestBody["text"][0])

		// empty slices will always have len() of 1 apparently, so check for empty string
		if cmdText[0] == "" {

			helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Missing Want ID"))
			return
		}

		if len(cmdText) >= 1 {

			id, err := strconv.Atoi(cmdText[0])
			if err != nil || id < 1 {
				helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Invalid Want ID"))
				return
			}

			if err := controllers.DeleteWant(id); err != nil {
				helpers.RespondWithError(w, helpers.ArgParser(http.StatusInternalServerError, err.Error()))
				return
			}

			helpers.RespondWithJSON(w, helpers.ArgParser(http.StatusOK, map[string]interface{}{"result": "success"}))
			return
		}
	}
	// TODO: maybe put an error here?
}

func executeTests(w http.ResponseWriter, r *http.Request) {

	test := `{"message": "yay tests"}`

	controllers.Tests()

	json.NewEncoder(w).Encode(test)
}

func main() {

	controllers.OpenDatabase()

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/get-wants", get)
	r.HandleFunc("/get-wants/{id}", get)

	r.HandleFunc("/create-want", put)

	r.HandleFunc("/update-want", post)
	r.HandleFunc("/update-want/{id}", post)

	r.HandleFunc("/delete-want", delete)
	r.HandleFunc("/delete-want/{id}", delete)

	// r.HandleFunc("/tests", executeTests)

	log.Fatal(http.ListenAndServe(":8080", r))
}
