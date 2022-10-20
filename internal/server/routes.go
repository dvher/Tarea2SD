package server

import (
	"log"
	"net/http"
    "math/rand"

    "tarea2sd/pkg/miembro"
    "tarea2sd/pkg/venta"
    "tarea2sd/pkg/strange"
	"tarea2sd/internal/producer"
	"github.com/gin-gonic/gin"
)

func registerMember(c *gin.Context) { 

    member := new(miembro.Miembro)

    err := c.BindJSON(member)

    if err != nil {
        panic(err)
    }

    prod, err := producer.NewProducer(BrokerList)

    if err != nil {
        log.Panic(err)
    }

    defer prod.Close()

    memberBytes, err := member.MarshalJSON()

    if err != nil {
        log.Panic(err)
    }

    part := int32(0)

    if member.Premium {
        part = int32(1)
    }

    _, _, err = prod.SendMessage("Membresia", part, memberBytes)

    if err != nil {
        log.Panic(err)
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func registerSale(c *gin.Context) {

    sale := new(venta.Venta)

    err := c.BindJSON(sale)

    if err != nil {
        panic(err)
    }

    prod, err := producer.NewProducer(BrokerList)

    if err != nil {
        log.Panic(err)
    }

    defer prod.Close()

    saleBytes, err := sale.MarshalJSON()

    if err != nil {
        log.Panic(err)
    }

    _, _, err = prod.SendMessage("Venta", rand.Int31n(2), saleBytes)

    if err != nil {
        log.Panic(err)
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func registerStrange(c *gin.Context) {

    stranger := new(strange.Strange)

    err := c.BindJSON(stranger)

    if err != nil {
        panic(err)
    }

    prod, err := producer.NewProducer(BrokerList)

    if err != nil {
        log.Panic(err)
    }

    defer prod.Close()

    strangerBytes, err := stranger.MarshalJSON()

    if err != nil {
        log.Panic(err)
    }

    _, _, err = prod.SendMessage("Extrano", rand.Int31n(2), strangerBytes)

    if err != nil {
        log.Panic(err)
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ping(c *gin.Context) {
    c.String(http.StatusOK, "pong")
}
