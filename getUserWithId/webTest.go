package main

import (
	"fmt"

	"database/sql"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/denisenkom/go-mssqldb"
)

const dbconnection = "server=antibuddies.co362eqfasab.us-east-2.rds.amazonaws.com;user id=antibuddies;password=WeberStudent1;port=1433"

type User struct {
	ID            string
	FirstName     string
	LastName      string
	Username      string
	IsAdmin       string
}

type request struct {
	ID string `json:"ID"`
}

type response struct {
	UserID        string `json:"ID"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Username      string `json:"username"`
	IsAdmin       string `json:"isAdmin"`
}

func GetUserInDB(request request) (User, error) {

	var user User

	db, err := sql.Open("mssql", dbconnection)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.Query("USE tempdb; GO")

	user.ID = request.ID

	err = db.QueryRow(`USE antibuddies; SELECT firstName, lastName, username, isAdmin FROM dbo.Users WHERE user_id = '`+request.ID+`';`).Scan(&user.FirstName, &user.LastName, &user.Username, &user.IsAdmin)
	if err != nil {
		panic(err)
	}

	return user, err
}

func handler(request request) ([]response, error) {

	user, err := GetUserInDB(request)

	arrayResponse := []response{}

	if err != nil {
		resp := response{
			UserID:        fmt.Sprintf("error in getting user"),
			FirstName:     fmt.Sprintf("error in getting user"),
			LastName:      fmt.Sprintf("error in getting user"),
			Username:      fmt.Sprintf("error in getting user"),
			IsAdmin:       fmt.Sprintf("error in getting user"),
		}
		arrayResponse = append(arrayResponse, resp)
		return arrayResponse, nil
	} else {
		resp := response{
			UserID:        fmt.Sprintf(user.ID),
			FirstName:     fmt.Sprintf(user.FirstName),
			LastName:      fmt.Sprintf(user.LastName),
			Username:      fmt.Sprintf(user.Username),
			IsAdmin:       fmt.Sprintf(user.IsAdmin),
		}
		arrayResponse = append(arrayResponse, resp)
		return arrayResponse, nil
	}
}

func main() {
	lambda.Start(handler)
}