package user

import (
	"encoding/json"
	"errors"
	"github.com/BentleyOph/go-serverless/pkg/validators"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// User is a struct that represents a user in the database.
type User struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var (
	ErrorFailedToGetUser     = "failed to get user"
	ErrorInvalidUserData     = "invalid user data"
	FailedToUnmarshallRecord = "failed to unmarshall record"
	InvalidEmail             = "invalid email"
	CouldNotMarshall         = "could not marshall"
	CouldNotDelete           = "could not delete"
	CouldNotPut              = "could not put"
	UserAlreadyExists        = "user already exists"
	UserDoesNotExist         = "user does not exist"
)

// CreateUser is a function that creates a new user in the database.
func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var user User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}
	if !validators.IsEmail(user.Email) {
		return nil, errors.New(InvalidEmail)
	}
	currentUser, _ := GetUser(user.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) > 0 {
		return nil, errors.New(UserAlreadyExists)
	}
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(CouldNotMarshall)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(CouldNotPut)
	}
	return &user, nil

}

// GetUsers is a function that retrieves users from the database.
func GetUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := dynaClient.Scan(input) // Get all users
	if err != nil {
		return nil, errors.New("failed to get users")
	}
	users := []User{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, errors.New(CouldNotMarshall)
	}
	return &users, nil

}

// GetUser is a function that retrieves a user from the database.
func GetUser(email string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToGetUser)
	}
	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New("failed to unmarshal")
	}

	return item, nil

}

// UpdateUser is a function that updates a user in the database.
func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}
	currentUser, _ := GetUser(u.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) > 0 {
		return nil, errors.New(UserDoesNotExist)
	}
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(CouldNotMarshall)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(CouldNotPut)
	}
	return &u, nil
}

// DeleteUser is a function that deletes a user from the database.
func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]
	input := dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	_ , err :=dynaClient.DeleteItem(&input)
	if err != nil {
		return errors.New(CouldNotDelete)
	}
	return nil
}
