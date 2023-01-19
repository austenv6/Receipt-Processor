package main

import (
	"fmt"
	//"encoding/json"
	"io/ioutil"
	"log"
)

func main() {
	fmt.Println("First Go Program")
	// take in a json receipt
	content, err := ioutil.ReadFile("./examples/simple-receipt.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	} else {
		fmt.Println("Successfully opened the file")
	}

	if content == nil {
		fmt.Println("bad")
	}
	// return a json object with code generated id and points
}