package server

import (
	"log"
    "os"
    "strconv"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

var BrokerList []string

func getBrokerList() {
    for i := 1; ; i++ {
        broker := os.Getenv("BROKER_" + strconv.Itoa(i))
        if broker == "" {
            break
        }
        BrokerList = append(BrokerList, broker)
    }
}

func New() (*gin.Engine) {

    getBrokerList()

    if len(BrokerList) == 0 {
        log.Fatal("No brokers found")
    }

    config := sarama.NewConfig()

    admin, err := sarama.NewClusterAdmin(BrokerList, config)

    if err != nil {
        log.Panic(err)
    }

    admin.CreateTopic("Ventas", &sarama.TopicDetail{
        NumPartitions: 2,
        ReplicationFactor: 1,
    }, false)

    admin.CreateTopic("Stock", &sarama.TopicDetail{
        NumPartitions: 2,
        ReplicationFactor: 1,
    }, false)

    admin.CreateTopic("Coordenadas", &sarama.TopicDetail{
        NumPartitions: 2,
        ReplicationFactor: 1,
    }, false)

    admin.CreateTopic("Membresias", &sarama.TopicDetail{
        NumPartitions: 2,
        ReplicationFactor: 1,
    }, false)

    admin.CreateTopic("Extrano", &sarama.TopicDetail{
        NumPartitions: 2,
        ReplicationFactor: 1,
    }, false)


    defer admin.Close()

	r := gin.Default()

	r.GET("/ping", ping)
	r.POST("/member", registerMember)
	r.POST("/sale", registerSale)
	r.POST("/strange", registerStrange)

	return r
}
