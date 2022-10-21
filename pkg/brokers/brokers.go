package brokers

import (
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var Brokers []string

func init() {
	getBrokerList()
}

func getBrokerList() {
	for i := 1; ; i++ {
		broker := os.Getenv("BROKER_NET_" + strconv.Itoa(i))
		if broker == "" {
			break
		}
		Brokers = append(Brokers, broker)
	}

	if len(Brokers) == 0 {
		log.Panic("Brokers not found")
	}
	log.Println(Brokers)
}
