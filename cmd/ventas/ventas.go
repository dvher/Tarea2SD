package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dvher/Tarea2SD/internal/consumer"
	"github.com/dvher/Tarea2SD/pkg/brokers"
	"github.com/dvher/Tarea2SD/pkg/venta"
	_ "github.com/joho/godotenv/autoload"
)

func getVentas() (sales []venta.Venta) {

	cons, err := consumer.NewConsumer(brokers.Brokers)

	if err != nil {
		log.Panic(err)
	}

	defer cons.Close()

	consume, err := cons.ConsumeFromBeginning("Ventas", 0)

	if err != nil {
		log.Panic(err)
	}

	defer consume.Close()

	for msg := range consume.Messages() {
		var sale venta.Venta
		err = json.Unmarshal(msg.Value, &sale)
		if err != nil {
			log.Println(err)
			continue
		}
		sales = append(sales, sale)
		if consumer.IsLastMessage(consume, msg) {
			break
		}
	}

	consume2, err := cons.ConsumeFromBeginning("Ventas", 1)

	if err != nil {
		log.Panic(err)
	}

	defer consume2.Close()

	for msg := range consume2.Messages() {
		var sale venta.Venta
		err = json.Unmarshal(msg.Value, &sale)
		if err != nil {
			log.Println(err)
			continue
		}
		sales = append(sales, sale)
		if consumer.IsLastMessage(consume2, msg) {
			break
		}
	}

	for _, sale := range sales {
		txt, err := sale.JSONIndent()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(txt)
	}

	return

}

func main() {

	ticker := time.NewTicker(24 * time.Hour)

	for {
		go getVentas()
		<-ticker.C
	}

}
