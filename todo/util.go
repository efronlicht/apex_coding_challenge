package todo

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type BadRequestError string

func (err BadRequestError) Error() string {
	return string(err)
}

func writeErr(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case nil:
		panic("called writeErr on nil error") // this shouldn't happen
	case BadRequestError:
		http.Error(w, string(err), http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", 500) // don't expose internal errors
	}
}

func writeOKResp(w http.ResponseWriter, jsonResp []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, jsonResp)
}

type key string

const (
	USER key = "DB_USER"
	NAME key = "DB_NAME"
)

func (k key) env() (string, error) {
	v, ok := os.LookupEnv(string(k))
	if !ok {
		return "", fmt.Errorf("missing environment key $%s", k)
	} else if v == "" {
		return "", fmt.Errorf("empty environment key $%s", k)
	}
	return v, nil
}
func OpenDB() *sql.DB {
	var user, name string
	var err error
	if user, err = USER.env(); err != nil {
		log.Fatal(err)
	} else if name, err = NAME.env(); err != nil {
		log.Fatal(err)
	}
	db, _ := sql.Open("postgres",
		fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, name))
	return db
}
