package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"time"
)

func CreateDynamoDBClient(region string) *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Fatalf("Got error initializing AWS: %s", err)
	}
	svc := dynamodb.New(sess)

	return svc
}

func createTable(tableName string, client *dynamodb.DynamoDB) {
	tableInput := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Location"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Time"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Location"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Time"),
				KeyType:       aws.String("RANGE"),
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}

	if _, err := client.CreateTable(tableInput); err != nil {
		log.Fatal("Creating table error: ", err.Error())
	}

	log.Printf("Table %s successfully created!", tableName)
}

func deleteTable(tableName string, client *dynamodb.DynamoDB) {
	listTableInput := &dynamodb.ListTablesInput{}

	listTableOutput, err := client.ListTables(listTableInput)
	if err != nil {
		log.Fatal("Listing table error: ", err.Error())
	}

	for _, existedTableName := range listTableOutput.TableNames {

		// Delete table if true
		if *existedTableName == tableName {
			log.Printf("Table %s already exists.", tableName)
			deleteTableInput := &dynamodb.DeleteTableInput{
				TableName: aws.String(tableName),
			}

			if _, err := client.DeleteTable(deleteTableInput); err != nil {
				log.Fatal("Deleting table error: ", err.Error())
			}
		}
	}
}

func SetUpTableAWS(tableName string, awsClient *dynamodb.DynamoDB) {
	deleteTable(tableName, awsClient)
	time1 := time.NewTimer(10 * time.Second)
	<-time1.C
	createTable(tableName, awsClient)
	timer2 := time.NewTimer(10 * time.Second)
	<-timer2.C
}

func PutItemInput(tableName string, body Data, awsClient *dynamodb.DynamoDB) {
	// Create the item to be added to DynamoDB
	var item Item
	item.Location = body.Request.Query
	item.Time = body.Location.Localtime
	item.Data = body

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new network item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = awsClient.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	fmt.Println("Successfully added to table " + tableName)
}
