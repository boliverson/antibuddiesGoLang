package main

import (
    "fmt"

	"database/sql"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/denisenkom/go-mssqldb"
)

// Database connection string
const dbconnection = "server=antibuddies.co362eqfasab.us-east-2.rds.amazonaws.com;user id=antibuddies;password=WeberStudent1;port=1433"

// This function requires no request parameters, but the template requires that this struct exist
type request struct {
    ItemID      string `json:"citemID"`
    Section        string `json:"section"`
    Question       string `json:"question"`
    Difficulty     string `json:"difficulty"`
    CorrectAnswer  string `json:"correctAnswer"`
    AnswerDesc     string `json:"answerDesc"`
    Answer1     string `json:"answer1"`
    Answer2     string `json:"answer2"`
    Answer3     string `json:"answer3"`
    Answer4     string `json:"answer4"`
    Anum1       string  `json:"num1"`
    Anum2       string  `json:"num2"`
    Anum3       string  `json:"num3"`
    Anum4       string  `json:"num4"`
}
// Response object property and associated JSON key
type response struct {
    PracticeQuestion string `json:"ID"`
}
// Opens a Database connection and inserts a practiceQuestion. Returns a string containing the new PracticeQuestion's ID
func CreatePracticeQuestionInDB(practiceQuestion request) string {
    
    var id string
    
    db, err := sql.Open("mssql", dbconnection)
    if err != nil {
        panic(err)
    }
    // Holds the connection open until the surrounding function has finished executing
    defer db.Close()
    
    db.Query("USE antibuddies; GO")
    
    // Execute db stored procedure to add a new practiceQuestion
    _, err = db.Query(`USE antibuddies; INSERT INTO PracticeQuestions (citem_id, section, question, qdifficulty, atype, aresponse) 
                        VALUES ('` + practiceQuestion.ItemID + `', '` + practiceQuestion.Section + `', '` + practiceQuestion.Question + `', '` + practiceQuestion.Difficulty + `', '` + practiceQuestion.CorrectAnswer + `', '` + practiceQuestion.AnswerDesc + `');`)
    if err != nil {
        return "Error in Adding PracticeQuestions: " + err.Error()
    }
    // Get the new PracticeQuestionID by matching on a unique combination of properties
    err = db.QueryRow(`USE antibuddies; SELECT question_id FROM dbo.PracticeQuestions WHERE section = '` + practiceQuestion.Section + `' AND question = '` + practiceQuestion.Question + `' AND qdifficulty = '` + practiceQuestion.Difficulty + `' AND atype = '` + practiceQuestion.CorrectAnswer + `' AND aresponse = '` + practiceQuestion.AnswerDesc + `';`).Scan(&id)
    if err != nil {
        return "PracticeQuestion created successfully but the ID request failed." + err.Error()
	}
	
	_, err = db.Query(`USE antibuddies; INSERT INTO PQAnswers (question_id, qanswer, anum) 
                        VALUES ('` + id + `', '` + practiceQuestion.Answer1 + `', '` + practiceQuestion.Anum1 + `'),
                            ('` + id + `', '` + practiceQuestion.Answer2 + `', '` + practiceQuestion.Anum2 + `'),
                            ('` + id + `', '` + practiceQuestion.Answer3 + `', '` + practiceQuestion.Anum3 + `'),
                            ('` + id + `', '` + practiceQuestion.Answer4 + `', '` + practiceQuestion.Anum4 + `');`)
    if err != nil {
        return "Error in Adding PracticeQuestionAnswers" + err.Error()
    }
    
    // Return new practiceQuestion's ID
    return id
}

// AWS Lambda template function that does all the work for the lambda execution
func handler(practiceQuestion request) (response, error) {
    
    PracticeQuestionID := CreatePracticeQuestionInDB(practiceQuestion)

    if PracticeQuestionID != "Error in Adding PracticeQuestionAnswers" {
		return response{
			PracticeQuestion: fmt.Sprintf(PracticeQuestionID),
		}, nil
	} else {
		return response{
			PracticeQuestion: fmt.Sprintf("Creation of Practice Question %s failed", PracticeQuestionID),
		}, nil
	}
}

// AWS lambda standard main function
func main() {
    lambda.Start(handler)
}
