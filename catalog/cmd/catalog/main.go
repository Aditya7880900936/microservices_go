package main

import (
	"log"
	"time"

	"github.com/Aditya7880900936/microservices_go/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}

	var r catalog.Repository // ✅ correct declaration

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("failed to connect to elasticsearch: %s", err)
			return err
		}
		return nil
	})
	defer r.Close()

	log.Println("Listening on port 8080...")

	s := catalog.NewService(r)
	err = catalog.ListenGRPC(s, 8080)
	if err != nil {
		log.Fatalf("failed to listen gRPC server: %s", err)
	}
}
