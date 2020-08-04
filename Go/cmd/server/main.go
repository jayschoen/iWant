package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

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
				helpers.RespondWithError(w, helpers.ItemFormatter(err.Error()))
				return
			}

			helpers.RespondWithJSON(w, helpers.ListFormatter(wants))
			return

		}

		if len(cmdText) >= 1 {

			id, err := strconv.Atoi(cmdText[0])
			if err != nil || id < 1 {
				helpers.RespondWithError(w, helpers.ItemFormatter("Invalid Want ID"))
				return
			}

			want, err := controllers.GetWantByID(id)
			if err != nil {
				helpers.RespondWithError(w, helpers.ItemFormatter(err.Error()))
				return
			}

			helpers.RespondWithJSON(w, helpers.PointerItemFormatter(want))
			return
		}
	}
	helpers.RespondWithError(w, helpers.ItemFormatter("Missing Request Body"))
	return
}

// TODO: NEEDS A DIALOG POPUP IN SLACK
// TODO:
func post(w http.ResponseWriter, r *http.Request) {

	/* 	vars := mux.Vars(r)
	   	params := r.URL.Query()

	   	paramID, hasID := vars["id"] */

	requestBody := helpers.ParseSlackPayload(r)

	if len(requestBody) > 0 {

		cmdText := helpers.ParseSlackPayloadText(requestBody["text"][0])

		// empty slices will always have len() of 1 apparently, so check for empty string
		if cmdText[0] == "" {
			helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Parameters Are Required"))
			return
		}

		if len(cmdText) >= 1 {

			// missing want ID

			// invalid want ID

			// missing other params

			// update or err

			helpers.RespondWithJSON(w, helpers.ArgParser(http.StatusOK, map[string]interface{}{"result": "success"}))
			return
		}

	}
	/* if !hasID {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Missing Want ID"))
		return
	}

	id, err := strconv.Atoi(paramID)
	if err != nil || id < 1 {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Invalid Want ID"))
		return
	} */

	/* if len(params) == 0 {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Parameters Are Required"))
		return
	} */

	/* status, statusExists := params["status"]
	wants, wantsExists := params["wants"]
	targetTime, targetTimeExists := params["targetTime"]

	if !statusExists && !wantsExists && !targetTimeExists {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusBadRequest, "Missing Parameter - Request Must Include At Least status Or wants Or targetTime"))
		return
	}

	if err := controllers.UpdateWant(id, status, wants, targetTime); err != nil {
		helpers.RespondWithError(w, helpers.ArgParser(http.StatusInternalServerError, err.Error()))
		return
	} */

	/* 	helpers.RespondWithJSON(w, helpers.ArgParser(http.StatusOK, map[string]interface{}{"result": "success"}))
	   	return */

}

//insert
func put2(w http.ResponseWriter, r *http.Request) {
	requestBody := helpers.ParseSlackPayload(r)

	if len(requestBody) > 0 {

		command := requestBody["command"][0]
		token := requestBody["token"][0]
		triggerID := requestBody["trigger_id"][0]

		slackExternalPost(token, triggerID, command)

	}

}

func put(w http.ResponseWriter, userInput UserInput) {

	// check if userInput is empty
	if userInput == (UserInput{}) {
		helpers.RespondWithError(w, helpers.ItemFormatter("Parameters Are Required"))
		return
	}

	slackID := userInput.SlackID
	status := userInput.Status
	wants := userInput.Wants
	targetTime := userInput.TargetTime

	if err := controllers.InsertWant(slackID, status, wants, helpers.ParseTimeString(targetTime)); err != nil {
		helpers.RespondWithError(w, helpers.ItemFormatter(err.Error()))
		return
	}

	//helpers.RespondWithJSON(w, helpers.ItemFormatter("Success"))

	w.WriteHeader(http.StatusOK)
	return
}

func delete(w http.ResponseWriter, r *http.Request) {

	requestBody := helpers.ParseSlackPayload(r)

	if len(requestBody) > 0 {

		cmdText := helpers.ParseSlackPayloadText(requestBody["text"][0])

		// empty slices will always have len() of 1 apparently, so check for empty string
		if cmdText[0] == "" {

			helpers.RespondWithError(w, helpers.ItemFormatter("Missing Want ID"))
			return
		}

		if len(cmdText) >= 1 {

			id, err := strconv.Atoi(cmdText[0])
			if err != nil || id < 1 {
				helpers.RespondWithError(w, helpers.ItemFormatter("Invalid Want ID"))
				return
			}

			if err := controllers.DeleteWant(id); err != nil {
				helpers.RespondWithError(w, helpers.ItemFormatter(err.Error()))
				return
			}

			helpers.RespondWithJSON(w, helpers.ItemFormatter("success"))
			return
		}
	}
	// TODO: maybe put an error here?
}

type UserInput struct {
	ActionID   string
	SlackID    string
	Status     string
	Wants      string
	TargetTime string
}

func captureUserInput(w http.ResponseWriter, r *http.Request) {

	fmt.Println("\n**************\ncaptureUserInput\n")

	requestBody := helpers.ParseSlackPayload(r)

	// fmt.Println(requestBody)

	payload := requestBody["payload"]

	fmt.Println("********** BEFORE MARSHAL")
	fmt.Println(payload[0])

	var vals map[string]interface{}
	err := json.Unmarshal([]byte(payload[0]), &vals)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("********** AFTER MARSHAL")
	fmt.Println(vals)

	actionID := vals["view"].(map[string]interface{})["callback_id"]

	user := vals["user"]

	values := vals["view"].(map[string]interface{})["state"].(map[string]interface{})["values"]

	fmt.Println("********** values after view -> state ->")
	fmt.Println(values)

	userID := user.(map[string]interface{})["id"]

	status := values.(map[string]interface{})["status"].(map[string]interface{})["status"].(map[string]interface{})["value"]

	wants := values.(map[string]interface{})["wants"].(map[string]interface{})["wants"].(map[string]interface{})["value"]

	targetDate := values.(map[string]interface{})["targetDate"].(map[string]interface{})["targetDate"].(map[string]interface{})["selected_date"]

	targetHour := values.(map[string]interface{})["targetHour"].(map[string]interface{})["targetHour"].(map[string]interface{})["selected_option"].(map[string]interface{})["value"]

	targetMinute := values.(map[string]interface{})["targetMinute"].(map[string]interface{})["targetMinute"].(map[string]interface{})["selected_option"].(map[string]interface{})["value"]

	fmt.Println(actionID, userID, status, wants, targetDate, targetHour, targetMinute)

	userInput := UserInput{
		actionID.(string),
		userID.(string), // maybe this should be user or username??
		status.(string),
		wants.(string),
		targetDate.(string) + "T" + targetHour.(string) + ":" + targetMinute.(string) + ":00.000Z",
	}

	if actionID == "create" {
		put(w, userInput)
		return
	}

	if actionID == "update" {

	}

	helpers.RespondWithError(w, helpers.ItemFormatter("Unspecified Action"))
	return

}

func slackExternalPost(token string, triggerID string, command string) {

	origination := command

	// TODO: PUT THIS IN A CONFIG ************************
	auth := ""
	// ***************************************************

	modalInfo := controllers.ConstructModalInfo(triggerID, origination)

	url := "https://slack.com/api/views.open"
	var bearer = "Bearer " + auth

	request, err := http.NewRequest("POST", url, strings.NewReader(modalInfo))

	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		fmt.Println("POST ERROR")
		fmt.Println(request)
		fmt.Println(err.Error())
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()

	fmt.Println("\n**********\nendOfexternalPost")
	fmt.Printf("%s\n", data)
}

func executeTests(w http.ResponseWriter, r *http.Request) {

	/* test := `{"message": "yay tests"}`

	controllers.Tests()

	json.NewEncoder(w).Encode(test) */

	slackExternalPost("12345", "abcd1234efgh567.X", "iwant-add2")
}

func main() {

	controllers.OpenDatabase()

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/get-wants", get)
	r.HandleFunc("/get-wants/{id}", get)

	//r.HandleFunc("/create-want", put)
	r.HandleFunc("/create-want2", put2)

	r.HandleFunc("/update-want", post)
	r.HandleFunc("/update-want/{id}", post)

	r.HandleFunc("/delete-want", delete)
	r.HandleFunc("/delete-want/{id}", delete)

	r.HandleFunc("/tests", executeTests)

	r.HandleFunc("/slack/capture", captureUserInput)

	log.Fatal(http.ListenAndServe(":8080", r))
}
