package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	controllers "github.com/jayschoen/iWant-backend/controllers"
)

type Blocks struct {
	Section []interface{} `json:"blocks"`
}
type Section struct {
	Type_ string      `json:"type"`
	Text  interface{} `json:"text"`
}
type FieldsSection struct {
	Type_  string    `json:"type"`
	Fields []Section `json:"fields"`
}

func RespondWithJSON(w http.ResponseWriter, payload interface{}) {

	response, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(response))

	w.Header().Set("Content-Type", "application/json")
	// slack **always** wants a response... so respond OK and pass the real response as json
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, message interface{}) {
	RespondWithJSON(w, message)
}

func ParseTimeString(str string) time.Time {
	// apparently takes the place of YYYY-MM-DD HH:MM:SS etc
	layout := "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}

	return t
}

func ParseSlackPayload(request *http.Request) url.Values {

	defer request.Body.Close()

	requestBody, err := ioutil.ReadAll(request.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	requestBodyString := string(requestBody)

	parsedBody, err := url.ParseQuery(requestBodyString)

	if err != nil {
		fmt.Println(err.Error())
	}

	return parsedBody
}

func ParseSlackPayloadText(text string) []string {
	return strings.Split(text, " ")
}

func ItemFormatter(data interface{}) Blocks {
	var blocks Blocks

	section := Section{
		Type_: "section",
		Text: Section{
			Type_: "mrkdwn",
			Text:  fmt.Sprintf("%+v", data),
		},
	}

	blocks.Section = append(blocks.Section, section)

	return blocks
}

func PointerItemFormatter(data *controllers.IWantRow) Blocks {
	var blocks Blocks

	var fields FieldsSection
	fields.Type_ = "section"

	fmt.Println(data)

	id := fmt.Sprint(data.Id)
	slackName := data.SlackName
	urgency := data.Urgency
	wants := data.Wants
	created := data.Created
	target := data.AppointmentTime

	values := [6]string{id, slackName, urgency, wants, created, target}

	headers := [6]string{"wantID", "slackName", "urgency", "wants", "created", "appointmentTime"}
	for key, val := range headers {

		tmp := fmt.Sprintf("*%v:* _%v_", val, values[key])

		field := Section{
			Type_: "mrkdwn",
			Text:  tmp,
		}
		fields.Fields = append(fields.Fields, field)
	}

	blocks.Section = append(blocks.Section, fields)

	return blocks
}

func ListFormatter(rawData []controllers.IWantRow) Blocks {
	var blocks Blocks

	var fields FieldsSection
	fields.Type_ = "section"

	divider := struct {
		Type_ string `json:"type"`
	}{
		Type_: "divider",
	}

	for index, data := range rawData {

		id := fmt.Sprint(data.Id)
		slackName := data.SlackName
		urgency := data.Urgency
		wants := data.Wants
		created := data.Created
		target := data.AppointmentTime

		values := [6]string{id, slackName, urgency, wants, created, target}

		headers := [6]string{"wantID", "slackName", "urgency", "wants", "created", "appointmentTime"}

		for key, val := range headers {

			tmp := fmt.Sprintf("*%v:* _%v_", val, values[key])

			field := Section{
				Type_: "mrkdwn",
				Text:  tmp,
			}
			fields.Fields = append(fields.Fields, field)

		}

		blocks.Section = append(blocks.Section, fields)

		if index < len(rawData)-1 {
			blocks.Section = append(blocks.Section, divider)
		}

		fields.Fields = nil
	}

	return blocks
}

func CheckAuthorization(userName string) bool {
	authorizedUsers := strings.Split(os.Getenv("APP_ADMIN_USERS"), ",")

	for _, authorizedUser := range authorizedUsers {
		if userName == authorizedUser {
			return true
		}
	}

	return false
}

type Users struct {
	Ok      bool `json:"ok"`
	Members []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
}

func getUserList() Users {

	auth := os.Getenv("SLACK_TOKEN")

	payload := fmt.Sprintf(`{"token": "%v"}`, auth)

	urlString := "https://slack.com/api/users.list"
	var bearer = "Bearer " + auth

	request, err := http.NewRequest("POST", urlString, strings.NewReader(payload))

	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err.Error())
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()

	var vals Users
	err = json.Unmarshal([]byte(data), &vals)
	if err != nil {
		fmt.Println(err)
	}

	return vals
}

type DirectMessageID struct {
	Ok      bool `json:"ok"`
	Channel struct {
		ID string `json"id"`
	}
}

func getDirectMessageID(userName string) DirectMessageID {

	users := getUserList()
	userID := ""
	for i := 0; i < len(users.Members); i++ {
		if users.Members[i].Name == userName {
			userID = users.Members[i].ID
		}
	}

	auth := os.Getenv("SLACK_TOKEN")

	payload := fmt.Sprintf(`{"token": "%v", "users": "%v"}`, auth, userID)

	url := "https://slack.com/api/conversations.open"
	var bearer = "Bearer " + auth

	request, err := http.NewRequest("POST", url, strings.NewReader(payload))

	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err.Error())
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()

	var vals DirectMessageID
	err = json.Unmarshal([]byte(data), &vals)
	if err != nil {
		fmt.Println(err)
	}

	return vals
}

func SendUpdateNotificationToUser(userName string, iWantID int) {

	directMessageID := getDirectMessageID(userName).Channel.ID

	message := fmt.Sprintf(`{
		"channel": "%v",
		"text": "Your iWant (#%v) has been updated."
	}`, directMessageID, iWantID)

	auth := os.Getenv("SLACK_TOKEN")

	url := "https://slack.com/api/chat.postMessage"
	var bearer = "Bearer " + auth

	request, err := http.NewRequest("POST", url, strings.NewReader(message))

	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()
}
