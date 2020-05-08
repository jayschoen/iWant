package controllers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"time"
	"math/rand"
)

func Tests() {

	openDatabase()

	fake_slack_id := testing_rand_num()

	created := time.Now()
	target_time := created.Add(time.Hour * 1)

	insert_want(fake_slack_id, "Wants", "cheese", target_time)

	update_want( UpdateParams{1, "Updated_3", "thing_3", target_time} )

	delete_want(2)

	get_by_id := get_want_by_id(1)
	fmt.Println(get_by_id)

	get_all := get_all_wants()
	fmt.Println(get_all)
}

func get_want_by_id (
	id int,
) iWant_Row {

	var(
		slack_id int
		status string
		wants string
		created string
		target_time string
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
		err := get.Scan(&id, &slack_id, &status, &wants, &created, &target_time)

		if err != nil {
			panic( err.Error() )
		}
	}

	return iWant_Row { id, slack_id, status, wants, created, target_time }
}

func get_all_wants() []iWant_Row {

	var(
		id int
		slack_id int
		status string
		wants string
		created string
		target_time string
	)

	get, err := db.Query(
		"SELECT * FROM whatsup",
	)

	if err != nil {
		panic( err.Error() )
	}

	defer get.Close()

	var rows []iWant_Row
	for get.Next() {
		err := get.Scan(&id, &slack_id, &status, &wants, &created, &target_time)

		if err != nil {
			panic( err.Error() )
		}

		rows = append(rows, iWant_Row {id, slack_id, status, wants, created, target_time})
	}

	return rows
}

func insert_want(
	slack_id int,
	status string,
	wants string,
	target_time time.Time,
) {

	created := time.Now()

	insert, err := db.Query(
		"INSERT INTO whatsup (slack_id, status, wants, created, target_time ) VALUES (?, ?, ?, ?, ?)",
		slack_id, status, wants, created, target_time,
	)

	if err != nil {
		panic( err.Error() )
	}

	defer insert.Close()
}

func update_want(
	params UpdateParams,
) {

	update, err := db.Query(
		"UPDATE whatsup SET status = ?, wants = ?, target_time = ? WHERE id = ?",
		params.STATUS, params.WANTS, params.TARGET_TIME, params.ID,
	)

	if err != nil {
		panic( err.Error() )
	}

	defer update.Close()
}

func delete_want(
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

func testing_rand_num() int {
	rand.Seed(time.Now().UnixNano())
    min := 1
	max := 10000
	
    return (rand.Intn(max - min + 1) + min)
}

type iWant_Row struct {
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
func openDatabase() {
	db, _ = sql.Open("mysql", MYSQL_CREDS)
}
