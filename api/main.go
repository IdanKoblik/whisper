package main

import (
	"whisper-api/endpoints"
)

func main() {
	router := endpoints.SetupRouter()	
	router.Run(":8080")
}
