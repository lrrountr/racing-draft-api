package tests

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"

	"github.com/lrrountr/racing-draft-api/internal/config"
	"github.com/lrrountr/racing-draft-api/internal/handler"
	"github.com/lrrountr/racing-draft-api/internal/model"
	shared "github.com/lrrountr/racing-draft-api/pkg/racing-draft"
)

var (
	ServerPort uint16
	DBName     string
	Bucket     string
	URLBase    string

	DBClient model.DBClient
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Lshortfile | log.LstdFlags)
	flag.Parse()
	if testing.Short() {
		log.Println("Skipping HTTP tests due to short mode...")
	}

	ServerPort = GetRandomPort()
	URLBase = fmt.Sprintf("http://localhost:%d", ServerPort)

	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading conf: %s", err)
	}

	DBName = randomString()
	config.APIBaseURL = URLBase
	config.DB.DBName = DBName

	os.Setenv("RACING_DRAFT_DB_DBNAME", config.DB.DBName)
	err = EnsureDatabase(config)
	if err != nil {
		log.Printf("Could not ensure database for tests: %s", err)
	}
	DBClient, err = model.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to init DBClient: %s", err.Error())
	}
	config.Server.Address = "127.0.0.1"
	config.Server.Port = ServerPort

	go func() {
		err := handler.StartServer(config)
		if err != nil {
			log.Fatalf("Failed to start server: %s", err)
		}
	}()

	addr := fmt.Sprintf("http://127.0.0.1:%d/", ServerPort)
	waitForServer(addr)

	setupTestRequirements(addr, config)

	out := m.Run()
	if out != 0 {
		os.Exit(out)
	}

	stat := DBClient.Stat()
	if stat.AcquiredConns() != 0 {
		log.Printf("After finishing model tests, there were %d outstanding connections", stat.AcquireCount())
		os.Exit(1)
	}

}

func GetRandomPort() uint16 {
	sock, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer sock.Close()
	return uint16(sock.Addr().(*net.TCPAddr).Port)
}

func randomString() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}

func EnsureDatabase(conf config.Config) error {
	err := model.EnsureDatabase(conf.DB)
	if err != nil {
		return err
	}

	return nil
}

func waitForServer(url string) {
	t := time.After(5 * time.Minute)
	fails := 0
	for {
		select {
		case <-t:
			log.Println("Server failed to start...")
			os.Exit(1)
		default:
		}

		err := healthCheck(url)
		if err != nil {
			if fails > 10 {
				log.Println(err)
			}
			fails += 1
		} else {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func healthCheck(url string) error {
	c := &http.Client{Timeout: time.Minute}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code during health check: %d", resp.StatusCode)
	}
	return nil
}

func setupTestRequirements(serverURL string, config config.Config) {
}

func getAdminClient(t *testing.T) shared.AdminClient {
	c, err := shared.NewClient(
		shared.WithBaseURL(URLBase),
	)
	assert.NilError(t, err)

	return c
}
