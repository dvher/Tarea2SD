package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/dvher/Tarea2SD/internal/consumer"
	"github.com/dvher/Tarea2SD/pkg/brokers"
	"github.com/dvher/Tarea2SD/pkg/coordinates"
)

var profugos []coordinates.Coordinates
var finished chan bool
var members map[string]chan coordinates.Coordinates

func printProfugos() {

	fmt.Println("\nRegistro de Carritos Profugos: ")
	for _, v := range profugos {
		fmt.Printf("\t %v ", v.Coords)
		if v.Miembro != "" {
			fmt.Printf("de %s \n", v.Miembro)
		} else {
			fmt.Println("")
		}
	}
}

func handleCoords(ch chan coordinates.Coordinates) {

	waitTime := 1 * time.Minute
	timer := time.NewTimer(waitTime)
	var LastCoord coordinates.Coordinates

RECIEVECOORDS:

	for {

		select {

		case <-timer.C:
			addProfugo(LastCoord)
			printProfugos()

			close(members[LastCoord.Miembro])
			delete(members, LastCoord.Miembro)
			timer.Stop()

			break RECIEVECOORDS

		case LastCoord = <-ch:
			timer.Reset(waitTime)
		}

	}

}

func addProfugo(c coordinates.Coordinates) {
	for i, v := range profugos {
		if v.Miembro == c.Miembro && c.Miembro != "" {
			profugos[i].Coords = c.Coords
			return
		}
	}
	profugos = append(profugos, c)

}

func getCoordinates() {

	members = make(map[string]chan coordinates.Coordinates)

	cons, err := consumer.NewConsumerGroup(brokers.Brokers, "Coordenadas", sarama.OffsetOldest)

	if err != nil {
		log.Panic(err)
	}

	defer cons.Close()

	var carrito coordinates.Coordinates

	ch := consumer.ConsumerHandler{
		Ready: make(chan bool),
		F: func(msg *sarama.ConsumerMessage) {
			carrito = coordinates.Coordinates{}
			err := json.Unmarshal(msg.Value, &carrito)

			if err != nil {
				log.Panic(err)
			}

			if carrito.Miembro == "" {
				addProfugo(carrito)
				printProfugos()
				return
			}

			if _, ok := members[carrito.Miembro]; !ok {
				members[carrito.Miembro] = make(chan coordinates.Coordinates)
				go handleCoords(members[carrito.Miembro])
			}

			members[carrito.Miembro] <- carrito
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := cons.Consume(ctx, []string{"Coordenadas"}, &ch); err != nil {
				log.Panic(err)
			}

			if ctx.Err() != nil {
				return
			}

			ch.Ready = make(chan bool)

		}
	}()

	<-ch.Ready

	for {
		<-ctx.Done()
		break
	}

	cancel()
	wg.Wait()

	finished <- true

	return

}

func main() {

	finished = make(chan bool)

	go getCoordinates()

	<-finished

}
