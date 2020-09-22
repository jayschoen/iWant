package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	controllers "github.com/jayschoen/iWant/controllers"
	helpers "github.com/jayschoen/iWant/helpers"
)

func get(w http.ResponseWriter, r *http.Request) {

	if helpers.AuthenticateRequest(r) != true {
		helpers.RespondWithError(w, helpers.ItemFormatter("Invalid Slack Request"))
		return
	}

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

			if helpers.CheckAuthorization(requestBody["user_name"][0]) == false {
				var filteredWants []controllers.IWantRow
				for _, want := range wants {
					if want.SlackName == requestBody["user_name"][0] {
						filteredWants = append(filteredWants, want)
					}
				}

				wants = filteredWants
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

			if helpers.CheckAuthorization(requestBody["user_name"][0]) == false {
				if want.SlackName == requestBody["user_name"][0] {
					helpers.RespondWithJSON(w, helpers.PointerItemFormatter(want))
					return
				} else {
					helpers.RespondWithError(w, helpers.ItemFormatter("Unauthorized"))
					return
				}

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
	urgency := userInput.Urgency
	wants := userInput.Wants
	appointmentTime := userInput.appointmentTime

	if err := controllers.UpdateWant(wantID, urgency, wants, helpers.ParseTimeString(appointmentTime)); err != nil {
		helpers.RespondWithError(w, helpers.ItemFormatter(err.Error()))
		return
	}

	want, _ := controllers.GetWantByID(wantID)

	helpers.SendUpdateNotificationToUser(want.SlackName, wantID)
	w.WriteHeader(http.StatusOK)
	return

}

func preparePutOrPost(w http.ResponseWriter, r *http.Request) {

	if helpers.AuthenticateRequest(r) != true {
		helpers.RespondWithError(w, helpers.ItemFormatter("Invalid Slack Request"))
		return
	}

	requestBody := helpers.ParseSlackPayload(r)

	if len(requestBody) > 0 {

		command := requestBody["command"][0]
		token := requestBody["token"][0]
		triggerID := requestBody["trigger_id"][0]

		cmdText := helpers.ParseSlackPayloadText(requestBody["text"][0])
		if command == "/iwant-update" {

			if helpers.CheckAuthorization(requestBody["user_name"][0]) != true {
				helpers.RespondWithError(w, helpers.ItemFormatter("Not Authorized"))
				return
			}

			if cmdText[0] != "" {
				id, err := strconv.Atoi(cmdText[0])

				if err != nil || id < 1 {
					helpers.RespondWithError(w, helpers.ItemFormatter("Invalid Want ID"))
					return
				}

				slackExternalPost(token, triggerID, command, cmdText[0])
				return
			}
		}

		slackExternalPost(token, triggerID, command, "")
		return

	}

}

func put(w http.ResponseWriter, userInput UserInput) {

	// check if userInput is empty
	if userInput == (UserInput{}) {
		helpers.RespondWithError(w, helpers.ItemFormatter("Parameters Are Required"))
		return
	}

	slackName := userInput.SlackName
	urgency := userInput.Urgency
	wants := userInput.Wants
	appointmentTime := "0000-00-00T00:00:00.000Z"

	if err := controllers.InsertWant(slackName, urgency, wants, helpers.ParseTimeString(appointmentTime)); err != nil {
		helpers.RespondWithError(w, helpers.ItemFormatter(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func delete(w http.ResponseWriter, r *http.Request) {

	if helpers.AuthenticateRequest(r) != true {
		helpers.RespondWithError(w, helpers.ItemFormatter("Invalid Slack Request"))
		return
	}

	requestBody := helpers.ParseSlackPayload(r)

	if helpers.CheckAuthorization(requestBody["user_name"][0]) != true {
		helpers.RespondWithError(w, helpers.ItemFormatter("Not Authorized"))
		return
	}

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
	helpers.RespondWithError(w, helpers.ItemFormatter("Missing Request Body"))
	return
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
				Urgency struct {
					Urgency struct {
						SelectedOption struct {
							Value string `json:"value"`
						} `json:"selected_option"`
					} `json:"urgency"`
				} `json:"urgency"`
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
	ActionID        string
	WantID          int
	SlackName       string
	Urgency         string
	Wants           string
	appointmentTime string
}

func captureUserInput(w http.ResponseWriter, r *http.Request) {

	if helpers.AuthenticateRequest(r) != true {
		helpers.RespondWithError(w, helpers.ItemFormatter("Invalid Slack Request"))
		return
	}

	requestBody := helpers.ParseSlackPayload(r)

	payload := requestBody["payload"]

	var vals SlackJSON
	err := json.Unmarshal([]byte(payload[0]), &vals)
	if err != nil {
		fmt.Println(err)
	}

	actionID := vals.View.CallbackID

	user := vals.User

	values := vals.View.State.Values

	userName := user.Name
	wantIDString := values.WantID.WantID.Value
	urgency := values.Urgency.Urgency.SelectedOption.Value
	wants := values.Wants.Wants.Value
	targetDate := values.TargetDate.TargetDate.SelectedDate
	targetHour := values.TargetHour.TargetHour.SelectedOption.Value
	targetMinute := values.TargetMinute.TargetMinute.SelectedOption.Value

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
		urgency,
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

func slackExternalPost(token string, triggerID string, command string, optionalWantID string) {

	origination := command

	auth := os.Getenv("SLACK_TOKEN")

	modalInfo := controllers.ConstructModalInfo(triggerID, origination, optionalWantID)

	url := "https://slack.com/api/views.open"
	var bearer = "Bearer " + auth

	request, err := http.NewRequest("POST", url, strings.NewReader(modalInfo))

	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()

}

func main() {

	controllers.OpenDatabase()

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/get-wants", get)

	r.HandleFunc("/create-want", preparePutOrPost)

	r.HandleFunc("/update-want", preparePutOrPost)

	r.HandleFunc("/delete-want", delete)

	r.HandleFunc("/slack/capture", captureUserInput)

	log.Fatal(http.ListenAndServe(":8080", r))
}
