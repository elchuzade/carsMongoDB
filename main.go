package main

// Note the folder name is carsMongoDB

import (
	"fmt"
	"log"
	"net/http"

	"carsMongoDB/router"
)

func main() {
	fmt.Println("Running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router.Router()))
}
