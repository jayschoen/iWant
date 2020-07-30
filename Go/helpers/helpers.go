package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	controllers "github.com/jayschoen/iWant-backend/controllers"
)

type Blocks struct {
	Section []Section `json:"blocks"`
}
type Section struct {
	Type_ string      `json:"type"`
	Text  interface{} `json:"text"`
}

func RespondWithJSON(w http.ResponseWriter, payload interface{}) {

	response, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
	}

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
		panic(err)
	}

	requestBodyString := string(requestBody)

	fmt.Println(requestBodyString)

	parsedBody, err := url.ParseQuery(requestBodyString)
	if err != nil {
		panic(err)
	}

	return parsedBody
}

func ParseSlackPayloadText(text string) []string {
	return strings.Split(text, " ")
}

func ArgParser(code int, textValue interface{}) string {
	fmt.Println(textValue)

	return fmt.Sprintf("%v: %v", code, textValue)
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

	section := Section{
		Type_: "section",
		Text: Section{
			Type_: "mrkdwn",
			Text:  fmt.Sprintf("%+v", *data),
		},
	}

	blocks.Section = append(blocks.Section, section)

	return blocks
}

func ListFormatter(data []controllers.IWantRow) Blocks {
	var blocks Blocks

	for _, data := range data {
		temp := Section{
			Type_: "section",
			Text: Section{
				Type_: "mrkdwn",
				Text:  fmt.Sprintf("%+v", data),
			},
		}

		blocks.Section = append(blocks.Section, temp)
	}

	return blocks
}
