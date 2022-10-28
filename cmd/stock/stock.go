package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/dvher/Tarea2SD/internal/consumer"
	"github.com/dvher/Tarea2SD/pkg/brokers"
	"github.com/dvher/Tarea2SD/pkg/venta"
)

var (
	sig           chan bool
	porReponer    [5]string
	idxPorReponer int = 0
)

func reponerStocks() {
	fmt.Println("Se le debe reponer stock a:")
	for _, v := range porReponer {
		fmt.Printf("\t%s\n", v)
	}
	idxPorReponer = 0
}

func consumeStock() {

	cg, err := consumer.NewConsumerGroup(brokers.Brokers, "stock", sarama.OffsetOldest)

	if err != nil {
		log.Panic(err)
	}

	defer cg.Close()

	ch := consumer.ConsumerHandler{
		Ready: make(chan bool),
		F: func(msg *sarama.ConsumerMessage) {
			var v venta.Venta

			err := json.Unmarshal(msg.Value, &v)

			if err != nil {
				log.Panic(err)
			}

			if v.Stock < 20 {
				porReponer[idxPorReponer] = v.Maestro
				idxPorReponer++
			}

			if idxPorReponer >= 5 {
				reponerStocks()
			}

		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {

			if err := cg.Consume(ctx, []string{"Stock"}, &ch); err != nil {
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
