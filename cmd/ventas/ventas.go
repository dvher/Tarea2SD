package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/dvher/Tarea2SD/internal/consumer"
	"github.com/dvher/Tarea2SD/pkg/brokers"
	"github.com/dvher/Tarea2SD/pkg/venta"
)

var (
	outputFile string
)

type Maestro struct {
	VentasTotales   int     `json:"ventas_totales"`
	PromedioVentas  float64 `json:"promedio_ventas"`
	ClientesTotales int     `json:"clientes_totales"`
}

func processVentas() {

	ventas := getVentas()

	ventas_maestro := make(map[string]Maestro)

	getVentasTotales(ventas, &ventas_maestro)
	getPromedioVentas(ventas, &ventas_maestro)
	getClientesTotales(ventas, &ventas_maestro)

	if outputFile != "" {
		writeToFile(ventas_maestro)
		return
	}

	fmt.Println(string(getString(ventas_maestro)))

}

func getVentasTotales(ventas map[string][]venta.Venta, maestros *map[string]Maestro) {

	for k := range ventas {
		total := 0

		for _, sale := range ventas[k] {
			total += sale.Cantidad
		}

		if maestro, ok := (*maestros)[k]; ok {
			maestro.VentasTotales = total
			(*maestros)[k] = maestro
		} else {
			maestro := Maestro{0, 0, 0}
			maestro.VentasTotales = total
			(*maestros)[k] = maestro
		}
	}

}

func getPromedioVentas(ventas map[string][]venta.Venta, maestros *map[string]Maestro) {

	for k := range ventas {
		total := float64(0)

		for _, sale := range ventas[k] {
			total += float64(sale.Cantidad)
		}

		if maestro, ok := (*maestros)[k]; ok {
			maestro.PromedioVentas = total / float64(len(ventas[k]))
			(*maestros)[k] = maestro
		} else {
			maestro := Maestro{0, 0, 0}
			maestro.PromedioVentas = total / float64(len(ventas[k]))
			(*maestros)[k] = maestro
		}
	}

}

func getClientesTotales(ventas map[string][]venta.Venta, maestros *map[string]Maestro) {

	for k := range ventas {
		clientes := make(map[string]bool)

		for _, sale := range ventas[k] {
			clientes[sale.Cliente] = true
		}

		if maestro, ok := (*maestros)[k]; ok {
			maestro.ClientesTotales = len(clientes)
			(*maestros)[k] = maestro
		} else {
			maestro := Maestro{0, 0, 0}
			maestro.ClientesTotales = len(clientes)
			(*maestros)[k] = maestro
		}
	}

}

func getString(maestros map[string]Maestro) []byte {

	txt, err := json.MarshalIndent(maestros, "", "\t")

	if err != nil {
		log.Panic(err)
	}

	return txt
}

func writeToFile(maestros map[string]Maestro) {

	f, err := os.Create(outputFile)

	if err != nil {
		log.Panic(err)
	}

	_, err = f.Write(getString(maestros))

	if err != nil {
		log.Panic(err)
	}

}

func getVentas() (sales map[string][]venta.Venta) {

	sales = make(map[string][]venta.Venta)

	cons, err := consumer.NewConsumerGroup(brokers.Brokers, "ventas", sarama.OffsetOldest)

	if err != nil {
		log.Panic(err)
	}

	defer cons.Close()

	ch := consumer.ConsumerHandler{
		Ready: make(chan bool),
		F: func(msg *sarama.ConsumerMessage) {
			var v venta.Venta
			err := json.Unmarshal(msg.Value, &v)

			if err != nil {
				log.Panic(err)
			}

			sales[v.Maestro] = append(sales[v.Maestro], v)
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	if err := cons.Consume(ctx, []string{"Ventas"}, &ch); err != nil {
		log.Panic(err)
	}

	cancel()

	return

}

func main() {

	keepRunning := flag.Bool("r", false, "Keep the program running every 24 hours")
	flag.StringVar(&outputFile, "o", "", "File to write sales to")
	flag.Parse()

	if !*keepRunning {
		processVentas()
		return
	}

	ticker := time.NewTicker(24 * time.Hour)

	for {
		go processVentas()
		<-ticker.C
	}

}
