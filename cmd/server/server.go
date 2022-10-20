package main

import (
	"flag"
	"log"
	"tarea2sd/internal/server"
    _ "github.com/joho/godotenv/autoload"
)


func main() {

    addr := flag.String("a", ":8000", "Port to run HTTP server")
    flag.Parse()

    router := server.New()

    log.Printf("Server running in port localhost%s\n", *addr)
    log.Fatal(router.Run(*addr))
    
}
