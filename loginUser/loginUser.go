package main

import (
	"database/sql"

	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/denisenkom/go-mssqldb"
)

const dbconnection = "server=antibuddies.co362eqfasab.us-east-2.rds.amazonaws.com;user id=antibuddies;password=WeberStudent1;port=1433"

type authentication struct {
	PassHash string
	UserID   string
	CanLogin bool
}

type Request struct {
	Username string `json:"username"`
	Pass  string `json:"password"`
}

type Response struct {
	CanLogin bool   `json:"response"`
	UserID   string `json:"ID"`
}

func checkInDB(request Request) authentication {

	var dbResponse authentication

	db, err := sql.Open("mssql", dbconnection)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.QueryRow(`USE antibuddies; SELECT password, user_id FROM dbo.Users WHERE username = '`+request.Username+`';`).Scan(&dbResponse.PassHash, &dbResponse.UserID)
	if err != nil {
		dbResponse.CanLogin = false
		return dbResponse
	}

	if dbResponse.PassHash == request.Pass {
		dbResponse.CanLogin = true
		return dbResponse
	}

	dbResponse.CanLogin = false
	return dbResponse
}

func Handler(request Request) (Response, error) {

	var dbResponse authentication

	dbResponse = checkInDB(request)

	if dbResponse.CanLogin {
		return Response{CanLogin: true, UserID: dbResponse.UserID}, nil
	}

	return Response{CanLogin: false, UserID: "0"}, nil
}

func main() {
	lambda.Start(Handler)
}