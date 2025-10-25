package main

import (
	"context"
	"whisper-api/config"
)

var ctx = context.Background()

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	router := SetupRouter(cfg)
	err = router.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}
