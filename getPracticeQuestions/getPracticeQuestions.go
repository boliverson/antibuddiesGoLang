package main

import (
	"database/sql"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/denisenkom/go-mssqldb"
)

// Database connection string
const dbconnection = "server=antibuddies.co362eqfasab.us-east-2.rds.amazonaws.com;user id=antibuddies;password=WeberStudent1;port=1433"

// Event properties and associated JSON keys that will make up the items in the returned JSON array
type PracticeQuestion struct {
	ID			   string `json:"id"`
	Section        string `json:"section"`
	Question       string `json:"question"`
	Difficulty     string `json:"difficulty"`
	CorrectAnswer  string `json:"correctAnswer"`
	AnswerDesc     string `json:"answerDesc"`
}



// This function requires no request parameters, but the template requires that this struct exist
type request struct {
	CourseID   string `json:"courseID"`
	Difficulty string `json:"difficulty"`
}

// Go slice will be converted to a JSON array labeled "PracticeQuestions"
type response struct {
	PracticeQuestions []*PracticeQuestion `json:"PracticeQuestions"`
}

// Workhorse method of the Lambda function
func GetPracticeQuestionInDB(request request) ([]*PracticeQuestion, error) {

	practiceQuestions := []*PracticeQuestion{}

	db, err := sql.Open("mssql", dbconnection)
	if err != nil {
		panic(err)
	}

	// Keeps database connection open until surrounding function has finished executing
	defer db.Close()

	db.Query("USE antibuddies; GO")

	// ORM returns the multi-row results of the query into a navigable variable called "rows"
	rows, err := db.Query(`USE antibuddies;
	SELECT q.question_id, q.section, q.question, q.qdifficulty, q.atype, q.aresponse
    FROM
		Courses c
		JOIN CourseItems ci
		ON ci.course_id = c.course_id
		JOIN PracticeQuestions q
		ON q.citem_id = ci.citem_id
    WHERE
        q.qdifficulty = '` + request.Difficulty + `' AND c.course_id = '` + request.CourseID + `';`)
	if err != nil {
		panic(err)
	}

	// Throws an error if no rows returned
	if !rows.Next() {
		panic("No rows returned. Diffuculty: " + request.Difficulty + ", Course ID: " + request.CourseID)
	}

	// For at most 100 rows, map all row columns to an event property and append that event to the Events slice
	for i := 0; i < 100; i++ {

		prac := new(PracticeQuestion)

		rows.Scan(&prac.ID, &prac.Section, &prac.Question, &prac.Difficulty, &prac.CorrectAnswer, &prac.AnswerDesc)

		practiceQuestions = append(practiceQuestions, prac)

		if !rows.Next() {
			return practiceQuestions, nil
		}
	}

	// Returns the events slice and any error
	return practiceQuestions, err
}

// AWS Lambda template function that does all the work for the lambda execution
func handler(request request) (response, error) {

	practiceQuestions, err := GetPracticeQuestionInDB(request)
	if err != nil {
		panic(err)
	}
	return response{
		// Maps the Go slice to the Request and converts to JSON
		PracticeQuestions: practiceQuestions,
	}, nil
}

// AWS lambda standard main function
func main() {
	lambda.Start(handler)
}