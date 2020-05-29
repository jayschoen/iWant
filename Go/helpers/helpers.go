package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
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
