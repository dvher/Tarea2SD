package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dvher/Tarea2SD/internal/consumer"
	"github.com/dvher/Tarea2SD/pkg/brokers"
	"github.com/dvher/Tarea2SD/pkg/venta"
	_ "github.com/joho/godotenv/autoload"
)

//map de ventas y genera maestro aleatorio

func getVentas() (sales []venta.Venta) {

	ctx, cancel := context.WithCancel(context.Background())

	Handler := &consumer.ConsumerHandler{
		Ready: make(chan bool),
	}
	cons, err := consumer.NewConsumerGroup(brokers.Brokers, "ventas")

	if err != nil {
		log.Panic(err)
	}

	defer cons.Close()
	//<-Handler.Ready

	err = cons.Consume(ctx, "ventas", Handler)
	if err != nil {
		log.Panic(err)

	}

	cancel()

	for _, sale := range sales {
		txt, err := sale.JSONIndent()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(string(txt))
	}

	return

}

func main() {
	//Sera pq se ejecuta cada 24hrs
	ticker := time.NewTicker(24 * time.Hour)

	//queda pegado
	for {
		go getVentas()
		<-ticker.C
	}

}
