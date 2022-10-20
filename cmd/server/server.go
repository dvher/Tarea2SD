package main

import (
	"flag"
	"log"
    "strconv"
    "os"
	"tarea2sd/internal/server"
    _ "github.com/joho/godotenv/autoload"
)

func getBrokerList() []string {
    var brokers []string
    for i := 1; ; i++ {
        broker := os.Getenv("BROKER_" + strconv.Itoa(i))
        if broker == "" {
            break
        }
        brokers = append(brokers, broker)
    }
    return brokers
}

func main() {

    addr := flag.String("a", ":8000", "Port to run HTTP server")
    flag.Parse()

    router := server.New(getBrokerList())

    log.Printf("Server running in port localhost%s\n", *addr)
    log.Fatal(router.Run(*addr))
    
}
