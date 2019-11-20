package main

import (
	"fmt"

	"database/sql"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/denisenkom/go-mssqldb"
)

// Database connection string
const dbconnection = "server=antibuddies.co362eqfasab.us-east-2.rds.amazonaws.com;user id=antibuddies;password=WeberStudent1;port=1433"

// Request object properties and associated JSON keys
type request struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Username	  string `json:"username"`
	PassHash      string `json:"pass"`
}

// Response object property and associated JSON key
type response struct {
	User string `json:"ID"`
}

// Opens a Database connection and inserts a user. Returns a string containing the new User's ID
func CreateUserInDB(user request) string {

	var id string

	db, err := sql.Open("mssql", dbconnection)
	if err != nil {
		panic(err)
	}

	// Holds the connection open until the surrounding function has finished executing
	defer db.Close()

	db.Query("USE antibuddies; GO")

	// Execute db stored procedure to add a new user
	_, err = db.Query(`USE antibuddies; INSERT INTO Users (firstName, lastName, username, password, isAdmin) 
													VALUES ('` + user.FirstName + `', '` + user.LastName + `', '` + user.Username + `', '` + user.PassHash + `', 0);`)
	if err != nil {
		return "Error in Adding Users"
	}

	// Get the new UserID by matching on a unique combination of properties
	err = db.QueryRow(`USE antibuddies; SELECT user_id FROM dbo.Users WHERE username = '` + user.Username + `' AND password = '` + user.PassHash + `';`).Scan(&id)
	if err != nil {
		return "User created successfully but the ID request failed."
	}

	// Return new user's ID
	return id
}

// AWS Lambda template function that does all the work for the lambda execution
func handler(user request) (response, error) {

	UserID := CreateUserInDB(user)

	// If CreateUserInDB returns without error, return the new UserID
	if UserID != "Error in Adding Users" {
		return response{
			User: fmt.Sprintf(UserID),
		}, nil
	} else {
		return response{
			User: fmt.Sprintf("Creation of User %s failed", user.Username),
		}, nil
	}
}

// AWS lambda standard main function
func main() {
	lambda.Start(handler)
}