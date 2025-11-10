package main

import (
	"log"
	"time"

	"github.com/Aditya7880900936/microservices_go/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	config := Config{}
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		e, err := account.NewPostgresRepository(config.DatabaseURL) // ✅ fixed here
		if err != nil {
			log.Println(err)
			return err
		}
		r = e
		return nil
	})
	defer r.Close()

	log.Println("Listening on port 8080...")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
