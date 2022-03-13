package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func currencyConversion(curr string, v float64) (float64, error) {
	url := "http://api.exchangeratesapi.io/v1/latest?access_key=15934be9b0a839cc7998faeb9e3babf9"
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	currencies := map[string]interface{}{}
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		return 0, err
	}
	c, ok := currencies["rates"].(map[string]interface{})[curr].(float64)
	if !ok {
		return 0, fmt.Errorf("can't find such currency. please use another one")
	}

	result := v / currencies["rates"].(map[string]interface{})["RUB"].(float64) * c

	return result, nil
}

func dbConnection() (*sql.DB, error) {
	envs, err := godotenv.Read(".env")
	if err != nil {
		return new(sql.DB), err
	}
	host := envs["DB_HOST"]
	port, _ := strconv.Atoi(envs["DB_PORT"])
	db_user := envs["DB_USER"]
	password := envs["DB_PASS"]
	dbname := envs["DB_NAME"]

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, db_user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return new(sql.DB), err
	}
	return db, nil
}

func getUser(id int) (User, error) {
	db, err := dbConnection()
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	var user User
	query := `SELECT "id", "balance" FROM "users" WHERE "id" = $1`
	rows := db.QueryRow(query, id)
	err = rows.Scan(&user.Id, &user.Balance)
	if err != nil {
		return User{}, fmt.Errorf("no user with such id")
	}
	return user, nil
}

func createUser(id int, balance float64) (User, error) {
	db, err := dbConnection()
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	query := `INSERT INTO "users" ("id", "balance") VALUES ($1, $2)`
	_, err = db.Exec(query, id, balance)
	if err != nil {
		return User{}, err
	}
	return User{Id: id, Balance: balance}, nil
}

func updateUserBalance(id int, balance float64) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE "users" SET "balance"=$1 WHERE "id"=$2`
	_, err = db.Exec(query, balance, id)
	if err != nil {
		return err
	}
	return nil
}

func updateHistory(id int, t string, amount float64) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	dt := time.Now()
	query := `INSERT INTO "history" ("user_id", "type", "amount", "datetime") VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(query, id, t, amount, dt.Format("01-02-2006 15:04:05"))
	if err != nil {
		return err
	}
	return nil
}

func getHistory(id, page int, order_field, order string) ([]History, error) {
	var historyPage []History
	var history History

	db, err := dbConnection()
	if err != nil {
		return []History{}, err
	}
	defer db.Close()

	limit := 10
	offset := (page - 1) * 10
	query := fmt.Sprintf(
		"SELECT \"user_id\", \"type\", \"amount\", \"datetime\" FROM \"history\" WHERE \"user_id\" = %v ORDER BY \"%v\" %v LIMIT %v OFFSET %v",
		id,
		order_field,
		order,
		limit,
		offset)
	rows, err := db.Query(query)
	if err != nil {
		return []History{}, err
	}
	for i := offset + 1; rows.Next(); i++ {
		err := rows.Scan(&history.UserId, &history.Type, &history.Amount, &history.Datetime)
		if err != nil {
			return []History{}, err
		}
		history.Type = strings.TrimSpace(history.Type)
		history.Idx = i
		historyPage = append(historyPage, history)
	}

	return historyPage, nil
}
