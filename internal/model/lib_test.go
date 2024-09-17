package model

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/lrrountr/racing-draft-api/internal/config"
)

var Client DBClient

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		log.Println("Skipping handler tests due to short mode...")
		os.Exit(0)
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Use a random DB for this test, so we know we are the only ones using it.
	config.DB.DBName = randomString()
	err = EnsureDatabase(config.DB)
	if err != nil {
		log.Printf("could not ensure database for tests: %s", err)
		os.Exit(1)
	}

	dbClient, err := NewClient(config)
	if err != nil {
		log.Printf("could not connect to database for tests: %s", err)
		os.Exit(1)
	}

	Client = dbClient

	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}

	stat := dbClient.Stat()
	if stat.AcquiredConns() != 0 {
		log.Printf("After finishing model tests, there were %d outstanding connections", stat.AcquireCount())
		os.Exit(1)
	}

	// After all our tests are down, test our model down + model up functions.
	// This also will catch issues where our final model test has a TX leak
	err = dbClient.Reinit()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func randomString() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}

func randBool() bool {
	return rand.Intn(2) == 1
}
