package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dvher/Tarea2SD/internal/consumer"
	"github.com/dvher/Tarea2SD/pkg/brokers"
	"github.com/dvher/Tarea2SD/pkg/venta"
	_ "github.com/joho/godotenv/autoload"
)

//map de ventas y genera maestro aleatorio

func getVentas() (sales []venta.Venta) {

	ctx, cancel := context.WithCancel(context.Background())

	handler := &consumer.ConsumerHandler{
		Ready: make(chan bool),
	}

	cons, err := consumer.NewConsumerGroup(brokers.Brokers, "ventas")

	if err != nil {
		log.Panic(err)
	}

	defer cons.Close()
	//<-Handler.Ready

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := cons.Consume(ctx, []string{"ventas"}, handler); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			handler.Ready = make(chan bool)
		}
	}()

	<-handler.Ready // Await till the consumer has been set up

	keepRunning := true

	for keepRunning {
		select {
		case <-ctx.Done():
			keepRunning = false
		}

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
