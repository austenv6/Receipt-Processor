# Receipt Processor

Receipt processor using golang for Fetch application.

This service makes use of github.com/google/uuid, so run this first > go get github.com/google/uuid

To run this processor, go the the receipts directory and enter either of the following into the command line:
For POSTs > go run process.go post <relative path of json file to process> 
  Example > go run process.go post ./examples/simple-receipt.json

For GETs > go run process.go get <UUID>
  Example > go run process.go get 041b70be-34a9-45fe-a155-25b33ae2b24f

Created by Austen Vizcarra