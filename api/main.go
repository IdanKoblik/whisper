package main

import (
	"fmt"
	"whisper-api/config"
	"whisper-api/endpoints"
)

func main() {
	cfg, err := config.ConfigReader{}.ReadConfig()
	if err != nil {
		fmt.Printf("\nCannot read config file: %v\n", err)
		return
	}

	router := endpoints.SetupRouter(&cfg)	
	router.Run(cfg.Addr)
}
