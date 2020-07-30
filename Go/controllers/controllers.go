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

	fakeSlackID := testingRandNum()

	created := time.Now()
	targetTime := created.Add(time.Hour * 1)

	InsertWant(fakeSlackID, "Wants", "cheese", targetTime)

	var status, wants, time []string
	status = append(status, "string_3")
	wants = append(wants, "thing_3")
	time = append(time, targetTime.String())

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
		slackID    int
		status     string
		wants      string
		created    string
		targetTime string
	)

	err := db.QueryRow(
		"SELECT * FROM whatsup WHERE id = ?",
		id,
	).Scan(&id, &slackID, &status, &wants, &created, &targetTime)

	if err != nil {
		return nil, err
	}

	return &IWantRow{id, slackID, status, wants, created, targetTime}, nil
}

func GetAllWants() ([]IWantRow, error) {

	var (
		id         int
		slackID    int
		status     string
		wants      string
		created    string
		targetTime string
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
		if err := get.Scan(&id, &slackID, &status, &wants, &created, &targetTime); err != nil {
			return nil, err
		}

		rows = append(rows, IWantRow{id, slackID, status, wants, created, targetTime})
	}

	return rows, nil
}

func InsertWant(
	slackID int,
	status string,
	wants string,
	targetTime time.Time,
) error {

	created := time.Now()

	insert, err := db.Query(
		"INSERT INTO whatsup (slackID, status, wants, created, targetTime ) VALUES (?, ?, ?, ?, ?)",
		slackID, status, wants, created, targetTime,
	)

	if err != nil {
		return err
	}

	defer insert.Close()

	return nil

}

func UpdateWant(
	id int,
	statusRaw []string,
	wantsRaw []string,
	targetTimeRaw []string,
) error {

	var status, wants, targetTime string
	if len(statusRaw) > 0 {
		status = statusRaw[0]
	}
	if len(wantsRaw) > 0 {
		wants = wantsRaw[0]
	}
	if len(targetTimeRaw) > 0 {
		targetTime = targetTimeRaw[0]
	}

	m := map[string]interface{}{"status": status, "wants": wants, "targetTime": targetTime}
	var values []interface{}
	var set []string
	for _, k := range []string{"status", "wants", "targetTime"} {
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
	Id         int
	SlackID    int
	Status     string
	Wants      string
	Created    string
	TargetTime string
}

// Global DB
var db *sql.DB

const mysqlCreds = "docker:docker@tcp(172.19.0.2:3306)/iWant_db"

func OpenDatabase() {
	db, _ = sql.Open("mysql", mysqlCreds)
}
