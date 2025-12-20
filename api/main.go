package main

import (
	"log"
	"whisper-api/communication"
	"whisper-api/config"
	"whisper-api/endpoints"

	"github.com/coreos/go-systemd/daemon"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	router := endpoints.SetupRouter(cfg)
	err = router.Run(cfg.Addr)
	if err != nil {
		if _, notifyErr := daemon.SdNotify(false, "STATUS=Service failed; ERRNO=1"); notifyErr != nil {
			log.Printf("Failed to notify systemd: %v\n", notifyErr)
		}

		panic(err)
	}

	go communication.HandleHeartbeat(cfg)

	if _, err := daemon.SdNotify(false, "READY=1"); err != nil {
		log.Printf("Failed to notify systemd: %v\n", err)
	} else {
		log.Println("Notified systemd that the service is READY=1")
	}
}
