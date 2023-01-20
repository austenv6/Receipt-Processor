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
)

type Item struct {
	ShortDescription string 
	Price string 
}

type Receipt struct {
	Retailer string 
	PurchaseDate string 
	PurchaseTime string 
	Total string 
	Items []Item 
}
func main() {
	if (len(os.Args) != 2) {
		fmt.Println("Usage: go run process.go <relative path.json>")
		os.Exit(1)
	}
	fmt.Println(os.Args[1])
	// take in a json receipt
	receiptFile, err := ioutil.ReadFile(os.Args[1]);
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		os.Exit(1)
	} 

	var payload Receipt
	err = json.Unmarshal(receiptFile, &payload)
	if err != nil {
		log.Fatal("Error parsing json", err)
		os.Exit(1)
	}

	fmt.Println("retailer: ", payload.Retailer)
	fmt.Println("purchase date: ", payload.PurchaseDate)
	fmt.Println("purchase time: ", payload.PurchaseTime)
	fmt.Println("total: ", payload.Total)
	for i := 0; i < len(payload.Items); i++ {
		fmt.Println("short description: ", payload.Items[i].ShortDescription)
		fmt.Println("price: ", payload.Items[i].Price)
	}

	// count alphanumeric characters in retailer name
	var pointCount = 0
	for i := 0; i < len(payload.Retailer); i++ {
		if ((payload.Retailer[i] > 47 && payload.Retailer[i] < 58) || (payload.Retailer[i] > 64 && payload.Retailer[i] < 123)) {
			pointCount++
		}
	}

	fmt.Println("count for retailer: ", pointCount)

	// give fifty points if total is dollar amt with no cents
	floatTotal, err := strconv.ParseFloat(payload.Total, 64)
	if err != nil {
		log.Fatal("Error converting total to float", err)
	}
	fmt.Println(math.Mod(floatTotal,10))
	if (math.Trunc(floatTotal) == floatTotal) {
		pointCount += 50;
	}

	fmt.Println("count for total + retailer: ", pointCount)

	// give 25 pts if total is a multiple of .25
	if (math.Mod(floatTotal,.25) == 0) {
		pointCount +=25
	}

	fmt.Println("count for total + retailer + .25: ", pointCount)

	// 5 pts for every two items on the receipt
	pointCount += ((len(payload.Items)/2) * 5)

	fmt.Println("count with num items: ", pointCount)

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

	fmt.Println("count with trim/3", pointCount)

	// 6 pts if day in purchase date is odd
	fmt.Println("purchase date: ", payload.PurchaseDate)
	fmt.Println("purchase time: ", payload.PurchaseTime)
	date, err := time.Parse("2006-01-02", payload.PurchaseDate)
	if err != nil {
		log.Fatal("Error converting string to date")
		os.Exit(1)
	}

	if date.Day() % 2 != 0 {
		pointCount += 6
	}

	fmt.Println("count with date is odd: ", pointCount)

	//

	// return a json object with code generated id and points
}