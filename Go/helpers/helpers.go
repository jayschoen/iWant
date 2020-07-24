package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RespondWithJSON(w http.ResponseWriter, payload interface{}) {

	wrapped_payload := map[string]interface{}{"response": payload}

	response, _ := json.Marshal(wrapped_payload)

	w.Header().Set("Content-Type", "application/json")
	// slack **always** wants a response... so respond OK and pass the real response as json
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, message interface{}) {
	RespondWithJSON(w, map[string]interface{}{"error": message})
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

func ArgParser(code int, textValue interface{}) interface{} {
	return map[string]interface{}{"code": code, "text": textValue}
}
