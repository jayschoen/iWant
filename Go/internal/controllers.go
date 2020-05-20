package controllers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"time"
	"math/rand"
)

func Tests() {

	fakeSlackId := testingRandNum()

	created := time.Now()
	targetTime := created.Add(time.Hour * 1)

	insertWant(fakeSlackId, "Wants", "cheese", targetTime)

	updateWant( UpdateParams{1, "Updated_3", "thing_3", targetTime} )

	deleteWant(2)

	getById := GetWantById(1)
	fmt.Println(getById)

	getAll := GetAllWants()
	fmt.Println(getAll)
}

func GetWantById (
	id int,
) iWantRow {
	
	var(
		slackId int
		status string
		wants string
		created string
		targetTime string
	)

	get, err := db.Query(
		"SELECT * FROM whatsup WHERE id = ?",
		id,
	)

	if err != nil {
		panic( err.Error() )
	}

	defer get.Close()

	for get.Next() {
		err := get.Scan(&id, &slackId, &status, &wants, &created, &targetTime)

		if err != nil {
			panic( err.Error() )
		}
	}

	return iWantRow { id, slackId, status, wants, created, targetTime }
}

func GetAllWants() []iWantRow {

	var(
		id int
		slackId int
		status string
		wants string
		created string
		targetTime string
	)

	get, err := db.Query(
		"SELECT * FROM whatsup",
	)

	if err != nil {
		panic( err.Error() )
	}

	defer get.Close()

	var rows []iWantRow
	for get.Next() {
		err := get.Scan(&id, &slackId, &status, &wants, &created, &targetTime)

		if err != nil {
			panic( err.Error() )
		}

		rows = append(rows, iWantRow {id, slackId, status, wants, created, targetTime})
	}

	return rows
}

func insertWant(
	slackId int,
	status string,
	wants string,
	targetTime time.Time,
) {

	created := time.Now()

	insert, err := db.Query(
		"INSERT INTO whatsup (slackId, status, wants, created, targetTime ) VALUES (?, ?, ?, ?, ?)",
		slackId, status, wants, created, targetTime,
	)

	if err != nil {
		panic( err.Error() )
	}

	defer insert.Close()
}

func updateWant(
	params UpdateParams,
) {

	update, err := db.Query(
		"UPDATE whatsup SET status = ?, wants = ?, targetTime = ? WHERE id = ?",
		params.STATUS, params.WANTS, params.TARGET_TIME, params.ID,
	)

	if err != nil {
		panic( err.Error() )
	}

	defer update.Close()
}

func deleteWant(
	id int,
) {

	delete, err := db.Query(
		"DELETE FROM whatsup WHERE id = ?",
		id,
	)

	if err != nil {
		panic( err.Error() )
	}

	defer delete.Close()
}

func testingRandNum() int {
	rand.Seed(time.Now().UnixNano())
    min := 1
	max := 10000
	
    return (rand.Intn(max - min + 1) + min)
}

type iWantRow struct {
	ID int
	SLACK_ID int
	STATUS string
	WANTS string
	CREATED string
	TARGET_TIME string
}

type UpdateParams struct {
	ID int
	STATUS string
	WANTS string
	TARGET_TIME time.Time
}

// Global DB
var db *sql.DB

const MYSQL_CREDS = "docker:docker@tcp(172.19.0.2:3306)/iWant_db"

// open DB connection globally
func OpenDatabase() {
	db, _ = sql.Open("mysql", MYSQL_CREDS)
}
