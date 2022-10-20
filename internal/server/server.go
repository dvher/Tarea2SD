package server

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)


func New(brokerList []string) (*gin.Engine) {

    config := sarama.NewConfig()
    

    admin, err := sarama.NewClusterAdmin(brokerList, config)

    if err != nil {
        log.Panic(err)
    }

    admin.CreateTopic("Ventas", &sarama.TopicDetail{
        NumPartitions: 2,
        ReplicationFactor: 0,
    }, false)

    admin.CreateTopic("Stock", &sarama.TopicDetail{
        NumPartitions: 2,
        ReplicationFactor: 0,
    }, false)

    admin.CreateTopic("Coordenadas", &sarama.TopicDetail{
        NumPartitions: 2,
        ReplicationFactor: 0,
    }, false)

    admin.CreateTopic("Membresias", &sarama.TopicDetail{
        NumPartitions: 2,
        ReplicationFactor: 0,
    }, false)

    defer admin.Close()

	r := gin.Default()

	r.GET("/ping", ping)
	r.POST("/member", registerMember)
	r.POST("/sale", registerSale)
	r.POST("/strange", registerStrange)

	return r
}
