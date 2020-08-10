package controllers

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"math/rand"
	"time"
)

// Tests TODO: remove this later
func Tests() {

	fakeSlackID := "aosidjfoasjd" //testingRandNum()

	created := time.Now()
	appointmentTime := created.Add(time.Hour * 1)

	InsertWant(fakeSlackID, "Wants", "cheese", appointmentTime)

	var status, wants string
	status = "string_3"     //append(status, "string_3")
	wants = "thing_3"       //append(wants, "thing_3")
	time := appointmentTime //append(time, appointmentTime.String())

	UpdateWant(1, status, wants, time)

	DeleteWant(2)

	getByID, err := GetWantByID(1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(getByID)

	getAll, err := GetAllWants()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(getAll)
}

func GetWantByID(
	id int,
) (*IWantRow, error) {

	var (
		slackName       string
		status          string
		wants           string
		created         string
		appointmentTime string
	)

	err := db.QueryRow(
		"SELECT * FROM whatsup WHERE id = ?",
		id,
	).Scan(&id, &slackName, &status, &wants, &created, &appointmentTime)

	if err != nil {
		return nil, err
	}

	return &IWantRow{id, slackName, status, wants, created, appointmentTime}, nil
}

func GetAllWants() ([]IWantRow, error) {

	var (
		id              int
		slackName       string
		status          string
		wants           string
		created         string
		appointmentTime string
	)

	get, err := db.Query(
		"SELECT * FROM whatsup",
	)

	if err != nil {
		return nil, err
	}

	defer get.Close()

	var rows []IWantRow
	for get.Next() {
		if err := get.Scan(&id, &slackName, &status, &wants, &created, &appointmentTime); err != nil {
			return nil, err
		}

		rows = append(rows, IWantRow{id, slackName, status, wants, created, appointmentTime})
	}

	return rows, nil
}

func InsertWant(
	slackName string,
	status string,
	wants string,
	appointmentTime time.Time,
) error {

	created := time.Now()

	insert, err := db.Query(
		"INSERT INTO whatsup (slackName, status, wants, created, appointmentTime ) VALUES (?, ?, ?, ?, ?)",
		slackName, status, wants, created, appointmentTime,
	)

	if err != nil {
		return err
	}

	defer insert.Close()

	return nil

}

func UpdateWant(
	id int,
	status string,
	wants string,
	appointmentTime time.Time,
) error {

	/* var status, wants, appointmentTime string
	if len(statusRaw) > 0 {
		status = statusRaw[0]
	}
	if len(wantsRaw) > 0 {
		wants = wantsRaw[0]
	}
	if len(appointmentTimeRaw) > 0 {
		appointmentTime = appointmentTimeRaw[0]
	} */

	m := map[string]interface{}{"status": status, "wants": wants, "appointmentTime": appointmentTime}
	var values []interface{}
	var set []string
	for _, k := range []string{"status", "wants", "appointmentTime"} {
		if v, ok := m[k]; ok {
			if v == "" {
				continue
			}
			values = append(values, v)
			set = append(set, fmt.Sprintf("%s = ?", k))
		}
	}

	values = append(values, id)
	query := "UPDATE whatsup SET " + strings.Join(set, ", ")
	query = query + " WHERE id = ?"
	update, err := db.Query(query, values...)

	if err != nil {
		return err
	}

	defer update.Close()

	return nil
}

func DeleteWant(
	id int,
) error {

	delete, err := db.Query(
		"DELETE FROM whatsup WHERE id = ?",
		id,
	)

	defer delete.Close()

	return err
}

func testingRandNum() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 10000

	return (rand.Intn(max-min+1) + min)
}

type IWantRow struct {
	Id              int
	SlackName       string
	Status          string
	Wants           string
	Created         string
	AppointmentTime string
}

// Global DB
var db *sql.DB

const mysqlCreds = "docker:docker@tcp(172.19.0.2:3306)/iWant_db"

func OpenDatabase() {
	db, _ = sql.Open("mysql", mysqlCreds)
}

func ConstructModalInfo(triggerID string, origination string) string {

	fmt.Println(origination)

	var callbackID, title, wantIDblock string = "", "", ""
	var optional bool

	if origination == "/iwant-add" {
		callbackID = "create"
		title = "Add iWant"
		optional = false
	}

	if origination == "/iwant-update" {
		callbackID = "update"
		title = "Update iWant"
		optional = true

		wantIDblock = `				
		{
			"block_id": "want_id",
			"type": "input",
			"optional": false,
			"element": {
				"type": "plain_text_input",
				"action_id": "want_id",
				"placeholder": {
					"type": "plain_text",
					"text": "WantID"
				}
			},
			"label": {
				"type": "plain_text",
				"text": "WantID"
			}
		},`

	}

	modalInfo := fmt.Sprintf(`{
		"trigger_id": "%[1]s",
		"view": {
			"title": {
				"type": "plain_text",
				"text": "%[2]s",
				"emoji": true
			},
			"submit": {
				"type": "plain_text",
				"text": "Submit",
				"emoji": true
			},
			"type": "modal",
			"callback_id": "%[3]s",
			"close": {
				"type": "plain_text",
				"text": "Cancel",
				"emoji": true
			},
			"blocks": [
				%s
				{
					"block_id": "status",
					"type": "input",
					"optional": %[7]v,
					"element": {
						"type": "plain_text_input",
						"action_id": "status",
						"placeholder": {
							"type": "plain_text",
							"text": "Status? (eg: Wants)"
						}
					},
					"label": {
						"type": "plain_text",
						"text": "Status"
					}
				},
				{
					"block_id": "wants",
					"type": "input",
					"optional": %[7]v,
					"element": {
						"type": "plain_text_input",
						"action_id": "wants",
						"placeholder": {
							"type": "plain_text",
							"text": "What does it want?!"
						}
					},
					"label": {
						"type": "plain_text",
						"text": "Wants"
					}
				},
				{
					"block_id": "targetDate",
					"type": "input",
					"optional": %[7]v,
					"element": {
						"type": "datepicker",
						"action_id": "targetDate"
					},
					"label": {
						"type": "plain_text",
						"text": "Select a date",
						"emoji": true
					}
				},
				%[5]v,
				%[6]v
			]
		}
	}`, triggerID, title, callbackID, wantIDblock, datepickerHour(optional), datepickerMinute(optional), optional)

	fmt.Println(modalInfo)

	return modalInfo

}

func datepickerHour(optional bool) string {

	template := `
	{
		"text": {
			"type": "plain_text",
			"text": "%v %v",
			"emoji": false
		},
		"value": "%v"
	},`

	hours := fmt.Sprintf(template, 6, "AM", 6)
	meridiem := "AM"
	var twelveHour int

	// 6am (including above) to 8pm?
	for i := 7; i < 20; i++ {

		if i == 12 {
			meridiem = "PM"
		}

		twelveHour = i
		if i > 12 {
			twelveHour = i - 12
		}

		hours = hours + fmt.Sprintf(template, twelveHour, meridiem, i)

	}

	hourSelect := fmt.Sprintf(`
	{
		"block_id": "targetHour",
		"type": "input",
		"optional": %v,
		"element": {
			"type": "static_select",
			"action_id": "targetHour",
			"placeholder": {
				"type": "plain_text",
				"text": "Select an hour",
				"emoji": true
			},
			"options": [
				%v
			]
		},
		"label": {
			"type": "plain_text",
			"text": "Select an hour",
			"emoji": true
		}
	}`, optional, hours)

	// fmt.Println(hourSelect)

	return hourSelect
}

func datepickerMinute(optional bool) string {

	minuteSelect := fmt.Sprintf(`
	{
		"block_id": "targetMinute",
		"type": "input",
		"optional": %v,
		"element": {
			"type": "static_select",
			"action_id": "targetMinute",
			"placeholder": {
				"type": "plain_text",
				"text": "Select a minute",
				"emoji": true
			},
			"options": [
				{
					"text": {
						"type": "plain_text",
						"text": "00",
						"emoji": false
					},
					"value": "00"
				},
				{
					"text": {
						"type": "plain_text",
						"text": "15",
						"emoji": false
					},
					"value": "15"
				},
				{
					"text": {
						"type": "plain_text",
						"text": "30",
						"emoji": false
					},
					"value": "30"
				},
				{
					"text": {
						"type": "plain_text",
						"text": "45",
						"emoji": false
					},
					"value": "45"
				},
			]
		},
		"label": {
			"type": "plain_text",
			"text": "Select a minute",
			"emoji": true
		}
	}`, optional)

	return minuteSelect
}
