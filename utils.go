package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func currencyConversion(curr string, v float64) (float64, error) {
	url := "http://api.exchangeratesapi.io/v1/latest?access_key=15934be9b0a839cc7998faeb9e3babf9"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err) // TODO обработать
	}

	currencies := map[string]interface{}{}
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		panic(err) // TODO обработать
	}
	c, ok := currencies["rates"].(map[string]interface{})[curr].(float64)
	if !ok {
		return 0, fmt.Errorf("can't find such currency. please use another one")
	}

	result := v / currencies["rates"].(map[string]interface{})["RUB"].(float64) * c

	return result, nil
}

func dbConnection() *sql.DB {
	envs, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env file") // TODO переменные окружения для БД
	}
	host := envs["DB_HOST"]
	port, _ := strconv.Atoi(envs["DB_PORT"])
	db_user := envs["DB_USER"]
	password := envs["DB_PASS"]
	dbname := envs["DB_NAME"]
	// log.Println(host, port, db_user, password, dbname)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, db_user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	return db
}

func getUser(id int) (User, error) {
	db := dbConnection()
	defer db.Close()
	var user User
	query := `SELECT "id", "balance" FROM "users" WHERE "id" = $1` // TODO id уникальное значение
	rows := db.QueryRow(query, id)
	err := rows.Scan(&user.Id, &user.Balance)
	if err != nil {
		// log.Println("Unable to execute the query. %v", err)
		return User{}, fmt.Errorf("no user with this id")
	}
	return user, nil
}

func createUser(id int, balance float64) (User, error) { // TODO доделать
	db := dbConnection()
	defer db.Close()
	query := `INSERT INTO "users" ("id", "balance") VALUES ($1, $2)`
	_, err := db.Exec(query, id, balance)
	if err != nil {
		// log.Println("Unable to execute the query. %v", err)
		return User{}, err
	}
	// log.Println("Inserted a single record %v", err)
	return User{Id: id, Balance: balance}, nil
}

func updateUserBalance(id int, balance float64) error {
	db := dbConnection()
	defer db.Close()

	_, err := getUser(id)
	if err != nil {
		return err
	}

	query := `UPDATE "users" SET "balance"=$1 WHERE "id"=$2`
	_, err = db.Exec(query, balance, id)
	if err != nil {
		log.Println("Unable to execute the query. %v", err)
		return err
	}
	log.Println("Updated a single record %v", err)
	return nil
}
