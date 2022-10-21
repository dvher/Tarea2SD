package server

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/dvher/Tarea2SD/pkg/brokers"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {

	config := sarama.NewConfig()

	admin, err := sarama.NewClusterAdmin(brokers.Brokers, config)

	if err != nil {
		log.Panic(err)
	}

	topics, err := admin.ListTopics()

	if err != nil {
		log.Panic(err)
	}

	if _, exists := topics["Ventas"]; !exists {
		err = admin.CreateTopic("Ventas", &sarama.TopicDetail{
			NumPartitions:     2,
			ReplicationFactor: 1,
		}, false)

		if err != nil {
			log.Panic(err)
		}

	}

	if _, exists := topics["Stock"]; !exists {
		err = admin.CreateTopic("Stock", &sarama.TopicDetail{
			NumPartitions:     2,
			ReplicationFactor: 1,
		}, false)

		if err != nil {
			log.Panic(err)
		}
	}

	if _, exists := topics["Coordenadas"]; !exists {
		err = admin.CreateTopic("Coordenadas", &sarama.TopicDetail{
			NumPartitions:     2,
			ReplicationFactor: 1,
		}, false)

		if err != nil {
			log.Panic(err)
		}
	}

	if _, exists := topics["Membresias"]; !exists {
		err = admin.CreateTopic("Membresias", &sarama.TopicDetail{
			NumPartitions:     2,
			ReplicationFactor: 1,
		}, false)

		if err != nil {
			log.Panic(err)
		}
	}

	defer admin.Close()

	r := gin.Default()

	r.GET("/ping", ping)
	r.POST("/member", registerMember)
	r.POST("/sale", registerSale)
	r.POST("/strange", registerStrange)

	return r
}
