package main

import (
	"context"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/dvher/Tarea2SD/internal/consumer"
	"github.com/dvher/Tarea2SD/pkg/brokers"
)

var sig chan bool

func consumeStock() {

	cg, err := consumer.NewConsumerGroup(brokers.Brokers, "stock", sarama.OffsetOldest)

	if err != nil {
		log.Panic(err)
	}

	defer cg.Close()

	ch := consumer.ConsumerHandler{
		Ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {

			if err := cg.Consume(ctx, []string{"Ventas"}, &ch); err != nil {
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

	sig <- true

}

func main() {

	sig = make(chan bool)

	go consumeStock()
	<-sig

}
