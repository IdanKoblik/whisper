package main

import (
	"whisper-api/config"
	"whisper-api/endpoints"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	router := endpoints.SetupRouter(cfg)
	err = router.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}
