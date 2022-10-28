package server

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/dvher/Tarea2SD/internal/producer"
	"github.com/dvher/Tarea2SD/pkg/brokers"
	"github.com/dvher/Tarea2SD/pkg/coordinates"
	"github.com/dvher/Tarea2SD/pkg/miembro"
	"github.com/dvher/Tarea2SD/pkg/venta"
	"github.com/gin-gonic/gin"
)

// maestro aleatorio
func registerMember(c *gin.Context) {

	member := new(miembro.Miembro)

	err := c.BindJSON(member)

	if err != nil {
		log.Panic(err)
		return
	}

	prod, err := producer.NewProducer(brokers.Brokers)

	if err != nil {
		log.Panic(err)
		return
	}

	defer prod.Close()

	memberBytes, err := member.JSON()

	if err != nil {
		log.Panic(err)
		return
	}

	part := int32(0)

	if member.Premium {
		part = int32(1)
	}

	_, _, err = prod.SendMessage("Membresias", part, memberBytes)

	if err != nil {
		log.Panic(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func registerSale(c *gin.Context) {

	sale := new(venta.Venta)

	err := c.BindJSON(&sale)

	if err != nil {
		log.Panic(err)
		return
	}

	prod, err := producer.NewProducer(brokers.Brokers)

	if err != nil {
		log.Panic(err)
		return
	}

	defer prod.Close()

	saleBytes, err := sale.JSON()

	if err != nil {
		log.Panic(err)
		return
	}

	_, _, err = prod.SendMessage("Ventas", rand.Int31n(2), saleBytes)

	if err != nil {
		log.Panic(err)
		return
	}

	_, _, err = prod.SendMessage("Stock", rand.Int31n(2), saleBytes)

	if err != nil {
		log.Panic(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func registerStrange(c *gin.Context) {

	coords := new(coordinates.Coordinates)

	err := c.BindJSON(coords)

	if err != nil {
		log.Panic(err)
		return
	}

	prod, err := producer.NewProducer(brokers.Brokers)

	if err != nil {
		log.Panic(err)
		return
	}

	defer prod.Close()

	strangerBytes, err := coords.JSON()

	if err != nil {
		log.Panic(err)
		return
	}

	_, _, err = prod.SendMessage("Coordenadas", 0, strangerBytes)

	if err != nil {
		log.Panic(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func registerCoords(c *gin.Context) {

	coords := new(coordinates.Coordinates)

	err := c.BindJSON(coords)

	if err != nil {
		log.Panic(err)
		return
	}

	prod, err := producer.NewProducer(brokers.Brokers)

	if err != nil {
		log.Panic(err)
		return
	}

	defer prod.Close()

	strangerBytes, err := coords.JSON()

	if err != nil {
		log.Panic(err)
		return
	}

	_, _, err = prod.SendMessage("Coordenadas", 1, strangerBytes)

	if err != nil {
		log.Panic(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
