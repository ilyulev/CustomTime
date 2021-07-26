package jsontime

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.mongodb.org/mongo-driver/bson"
)

func Test(t *testing.T) {
	loc, _ := time.LoadLocation("Pacific/Auckland")
	time.Local = loc
	tests := []struct {
		input        string
		expectedRFC  time.Time
		expectedCust time.Time
		expectedSQL  time.Time
	}{
		{
			input:        `{"t1":"2020-09-21T20:26:52.8585767","t2":"2020-09-15T14:45:33.3034643","t3":"2020-09-15T14:45:33Z"}`,
			expectedRFC:  time.Date(2020, time.September, 21, 20, 26, 52, 858576700, loc),
			expectedCust: time.Date(2020, time.September, 15, 14, 45, 33, 303464300, loc),
			expectedSQL:  time.Date(2020, time.September, 15, 14, 45, 33, 0, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var e struct {
				T1 JSONTime `json:"t1"`
				T2 JSONTime `json:"t2"`
				T3 JSONTime `json:"t3"`
			}
			err := json.Unmarshal([]byte(test.input), &e)
			if err != nil {
				t.Errorf("unexpected error during unmarshalling: %v", err)
			}
			if !test.expectedRFC.Equal(e.T1.Time) {
				t.Errorf("expected RFC3339 of %v, got %v", test.expectedRFC, e.T1)
			}
			if !test.expectedCust.Equal(time.Time(e.T2.Time)) {
				t.Errorf("expected \"2006-01-02T15:04:05.0000000\" of %v, got %v", test.expectedCust, e.T2)
			}
			if !test.expectedSQL.Equal(time.Time(e.T3.Time)) {
				t.Errorf("expected \"2006-01-02T15:04:05\" of %v, got %v", test.expectedSQL, e.T3)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	loc, _ := time.LoadLocation("Pacific/Auckland")
	time.Local = loc
	av := &dynamodb.AttributeValue{S: aws.String("2020-09-15T14:45:33")}

	actual := JSONTime{}
	expect := time.Date(2020, time.September, 15, 14, 45, 33, 0, time.UTC)
	err := dynamodbattribute.Unmarshal(av, &actual)
	fmt.Println(err, reflect.DeepEqual(expect, actual))

	// Output:
	// <nil> true
}

func TestMarshalBSON(t *testing.T) {
	loc, _ := time.LoadLocation("Pacific/Auckland")
	time.Local = loc

	//	actual := JSONTime{}
	expect := time.Date(2020, time.September, 15, 14, 45, 33, 0, time.UTC)
	av := TestBSON{T: JSONTime{expect}}
	//av := TestBSON{T: expect}
	v, err := bson.Marshal(av)
	fmt.Println(err, string(v))

	// Output:
	// <nil> true
}

type TestBSON struct {
	T JSONTime `bson:"t"`
	//T time.Time `bson:"t"`
}
