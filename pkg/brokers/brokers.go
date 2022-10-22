package brokers

import (
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

var Brokers []string

func init() {
	getBrokerList()
}

func getBrokerList() {
	broker := os.Getenv("BROKER_NET")

	Brokers = strings.Split(broker, ",")

	if len(Brokers) == 0 {
		log.Panic("Brokers not found")
	}
}
