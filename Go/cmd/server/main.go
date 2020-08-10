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

func post(w http.ResponseWriter, userInput UserInput) {

	// check if userInput is empty
	if userInput == (UserInput{}) {
		helpers.RespondWithError(w, helpers.ItemFormatter("Parameters Are Required"))
		return
	}

	wantID := userInput.WantID
	status := userInput.Status
	wants := userInput.Wants
	targetTime := userInput.TargetTime

	if err := controllers.UpdateWant(wantID, status, wants, helpers.ParseTimeString(targetTime)); err != nil {
		helpers.RespondWithError(w, helpers.ItemFormatter(err.Error()))
		return
	}

	//helpers.RespondWithJSON(w, helpers.ItemFormatter("Success"))

	w.WriteHeader(http.StatusOK)
	return

}

func preparePutOrPost(w http.ResponseWriter, r *http.Request) {
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

	slackName := userInput.SlackName
	status := userInput.Status
	wants := userInput.Wants
	targetTime := userInput.TargetTime

	if err := controllers.InsertWant(slackName, status, wants, helpers.ParseTimeString(targetTime)); err != nil {
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

type SlackJSON struct {
	User struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"user"`
	View struct {
		CallbackID string `json:"callback_id"`
		State      struct {
			Values struct {
				WantID struct {
					WantID struct {
						Value string `json:"value"`
					} `json:"want_id"`
				} `json:"want_id"`
				Status struct {
					Status struct {
						Value string `json:"value"`
					} `json:"status"`
				} `json:"status"`
				Wants struct {
					Wants struct {
						Value string `json:"value"`
					} `json:"wants"`
				} `json:"wants"`
				TargetDate struct {
					TargetDate struct {
						SelectedDate string `json:"selected_date"`
					} `json:"targetDate"`
				} `json:"targetDate"`
				TargetHour struct {
					TargetHour struct {
						SelectedOption struct {
							Value string `json:"value"`
						} `json:"selected_option"`
					} `json:"targetHour"`
				} `json:"targetHour"`
				TargetMinute struct {
					TargetMinute struct {
						SelectedOption struct {
							Value string `json:"value"`
						} `json:"selected_option"`
					} `json:"targetMinute"`
				} `json:"targetMinute"`
			} `json:"values"`
		} `json:"state"`
	} `json:"view"`
}

type UserInput struct {
	ActionID   string
	WantID     int
	SlackName  string
	Status     string
	Wants      string
	TargetTime string
}

func captureUserInput(w http.ResponseWriter, r *http.Request) {

	requestBody := helpers.ParseSlackPayload(r)

	payload := requestBody["payload"]

	fmt.Println("********** BEFORE MARSHAL")
	fmt.Println(payload[0])

	var vals SlackJSON
	err := json.Unmarshal([]byte(payload[0]), &vals)
	if err != nil {
		fmt.Println(err)
	}

	actionID := vals.View.CallbackID

	user := vals.User

	values := vals.View.State.Values

	userID := user.ID
	userName := user.Name
	wantIDString := values.WantID.WantID.Value
	status := values.Status.Status.Value
	wants := values.Wants.Wants.Value
	targetDate := values.TargetDate.TargetDate.SelectedDate
	targetHour := values.TargetHour.TargetHour.SelectedOption.Value
	targetMinute := values.TargetMinute.TargetMinute.SelectedOption.Value

	fmt.Println(actionID, wantIDString, userName, userID, status, wants, targetDate, targetHour, targetMinute)

	wantIDInt := 0
	if actionID == "update" {
		wantIDInt, err = strconv.Atoi(wantIDString)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	userInput := UserInput{
		actionID,
		wantIDInt,
		userName,
		status,
		wants,
		targetDate + "T" + targetHour + ":" + targetMinute + ":00.000Z",
	}

	if actionID == "create" {
		put(w, userInput)
		return
	}

	if actionID == "update" {
		post(w, userInput)
		return
	}

	helpers.RespondWithError(w, helpers.ItemFormatter("Unspecified Action"))
	return

}

func slackExternalPost(token string, triggerID string, command string) {

	origination := command

	// TODO: PUT THIS IN A CONFIG ************************
	auth := "xoxb-1227014263283-1212095310967-H9FQ4FPq8E3QCRMAHP1Wlhbf"
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

	//fmt.Println("\n**********\nendOfexternalPost")
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

	r.HandleFunc("/create-want", preparePutOrPost)

	r.HandleFunc("/update-want", preparePutOrPost)

	r.HandleFunc("/delete-want", delete)

	r.HandleFunc("/tests", executeTests)

	r.HandleFunc("/slack/capture", captureUserInput)

	log.Fatal(http.ListenAndServe(":8080", r))
}
