package main

import (
	"database/sql"
	
	"github.com/aws/aws-lambda-go/lambda"
    _ "github.com/denisenkom/go-mssqldb"
)

// Database connection string
const dbconnection = "server=antibuddies.co362eqfasab.us-east-2.rds.amazonaws.com;user id=antibuddies;password=WeberStudent1;port=1433"

// Event properties and associated JSON keys that will make up the items in the returned JSON array
type PracticeQuestionAnswer struct {
    ID             string `json:"id"`
    CorrectAnswer  string `json:"correctAnswer"`
    AnswerNum      string `json:"answerNum"`
    QuestionID     string `json:"questionID"`
}
// This function requires no request parameters, but the template requires that this struct exist
type request struct {
    QuestionID  string `json:"questionID"`
}
// Go slice will be converted to a JSON array labeled "PracticeQuestionAnswers"
type response struct {
    PracticeQuestionAnswers []*PracticeQuestionAnswer `json:"PracticeQuestionAnswers"`
}

// Workhorse method of the Lambda function
func GetPracticeQuestionAnswerInDB(request request) ([]*PracticeQuestionAnswer, error) {
	
	practiceQuestionAnswers := []*PracticeQuestionAnswer{}
	
	db, err := sql.Open("mssql", dbconnection)
    if err != nil {
        panic(err)
    }
// Keeps database connection open until surrounding function has finished executing
	defer db.Close()
	
	db.Query("USE antibuddies; GO")
	
	// ORM returns the multi-row results of the query into a navigable variable called "rows"
    rows, err := db.Query(`USE antibuddies;
	SELECT
		qanswer_id, qanswer, anum, question_id
    FROM
        PQAnswers
    WHERE
        question_id ='` + request.QuestionID + `';`)
    if err != nil {
        panic(err)
    }
    // Throws an error if no rows returned
    if !rows.Next() {
        panic("No rows returned. QuestionID: " + request.QuestionID)
    }
    // For at most 100 rows, map all row columns to an event property and append that event to the Events slice
    for i := 0; i < 100; i++ {
		
		answ := new(PracticeQuestionAnswer)
		
		rows.Scan(&answ.ID, &answ.CorrectAnswer, &answ.AnswerNum, &answ.QuestionID)
		
		practiceQuestionAnswers = append(practiceQuestionAnswers, answ)
		
		if !rows.Next() {
            return practiceQuestionAnswers, nil
        }
    }
    
    // Returns the events slice and any error
    return practiceQuestionAnswers, err
}
// AWS Lambda template function that does all the work for the lambda execution
func handler(request request) (response, error) {
	
	practiceQuestionAnswers, err := GetPracticeQuestionAnswerInDB(request)
    if err != nil {
        panic(err)
    }
    return response{
        // Maps the Go slice to the Request and converts to JSON
        PracticeQuestionAnswers: practiceQuestionAnswers,
    }, nil
}

// AWS lambda standard main function
func main() {
    lambda.Start(handler)
}