// package to process json receipts and handle get and post requests
package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"math"
	"strings"
	"time"
	"github.com/google/uuid"
)

// struct for items inside Receipt struct
type Item struct {
	ShortDescription string 
	Price string 
}

// struct to store incoming json objects
type Receipt struct {
	Retailer string 
	PurchaseDate string 
	PurchaseTime string 
	Total string 
	Items []Item 
}

// Points struct to output to json file after points calculated
type Points struct {
	Points string `json:"points"` // add field tag to get lowercase "points"
}

// UUID struct to eventually output UUID in json format
type UUID struct {
	UUID string `json:"id"`
}

//go run process.go post ./examples/simple-receipt.json
//go run process.go get <UUID>
//go get github.com/google/uuid

// POST handler function
func PerformPost() {
	// take in a json receipt
	receiptFile, err := ioutil.ReadFile(os.Args[2]);
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return
	} 

	var payload Receipt
	err = json.Unmarshal(receiptFile, &payload)
	if err != nil {
		log.Fatal("Error parsing json", err)
		return
	}

	// count alphanumeric characters in retailer name
	var pointCount = 0
	for i := 0; i < len(payload.Retailer); i++ {
		if ((payload.Retailer[i] > 47 && payload.Retailer[i] < 58) || (payload.Retailer[i] > 64 && payload.Retailer[i] < 123)) {
			pointCount++
		}
	}

	// give fifty points if total is dollar amt with no cents
	floatTotal, err := strconv.ParseFloat(payload.Total, 64)
	if err != nil {
		log.Fatal("Error converting total to float", err)
	}
	if (math.Trunc(floatTotal) == floatTotal) {
		pointCount += 50;
	}

	// give 25 pts if total is a multiple of .25
	if (math.Mod(floatTotal,.25) == 0) {
		pointCount +=25
	}

	// 5 pts for every two items on the receipt
	pointCount += ((len(payload.Items)/2) * 5)

	// if trimmed len of item descrip is mult 3, multiply by .2 and round up to nearest int
	for i := 0; i < len(payload.Items); i++ {
		strRes := strings.TrimSpace(payload.Items[i].ShortDescription)
		strLen := len(strRes)
		if (strLen % 3 == 0) {
			priceFloat, err := strconv.ParseFloat(payload.Items[i].Price, 64)
			if err != nil {
				log.Fatal("Error converting price to float", err)
				os.Exit(1)
			}
			pointCount += int(math.Round((priceFloat * 0.2) + 0.49)) // add .49 to ensure we round up, not down
		}
	}

	// 6 pts if day in purchase date is odd
	date, err := time.Parse("2006-01-02 15:04", payload.PurchaseDate + " " + payload.PurchaseTime)
	if err != nil {
		log.Fatal("Error converting string to date")
		return
	}

	if date.Day() % 2 != 0 {
		pointCount += 6
	}

	// if time of purchase is between 2pm and 4pm (200 to 359)
	if (date.Hour() >= 14 && date.Hour() < 16) {
		pointCount += 10
	}

	newUUID := (uuid.New()).String()
	
	// create points object and convert int to string
	pointsOut := &Points{
		Points: strconv.Itoa(pointCount),
	} 
	// Marshal pointsOut to json object
	output, err := json.Marshal(pointsOut)
	if err != nil {
		log.Fatal("Error creating JSON output object")
		return
	}
	// create UUID object
	uuidObj := &UUID{
		UUID: newUUID,
	}
	// Marshal uuidObj to json object
	uuidOutput, err := json.Marshal(uuidObj)
	if err != nil {
		log.Fatal("Error creating JSON UUID object")
		return
	}
	// create UUID directory
	err = os.Mkdir(newUUID, 0777)
	if err != nil {
		log.Fatal("Error creating UUID directory")
		return
	}
	// create points file under UUID directory
	err = ioutil.WriteFile(newUUID + "/points.json", output, 0644)
	if err != nil {
		log.Fatal("Error creating points file")
		return
	}
	// output UUID so user can grab it for get requests
	fmt.Println(string(uuidOutput))
	return
}

// get request handler
func PerformGet() {
	// grab the UUID passed in through command line and use it to try to locate corresponding points file
	checkUUID := strings.ToLower(os.Args[2])
	pointsFile, err := ioutil.ReadFile("./" + checkUUID + "/points.json")
	if err != nil {
		log.Fatal("Error when opening points file: ", err)
		return
	} 
	// unmarshal the points file
	var payload Points
	err = json.Unmarshal(pointsFile, &payload)
	if err != nil {
		log.Fatal("Error parsing json points file", err)
		return
	}
	// print points to console in json format
	fmt.Println(string(pointsFile))
	return
}

func main() {
	// verify correct number of args passed in
	if (len(os.Args) != 3) {
		fmt.Println("Usage: \nPost: go run process.go <relative path.json> \nGet: go run process.go get <UUID>")
		return
	}
	// determine what kind of request was made and call corresponding handler
	requestType := strings.ToUpper(os.Args[1])
	if requestType == "POST" {
		PerformPost()
	} else if requestType == "GET" {
		PerformGet()
	}
}