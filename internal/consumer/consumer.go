package consumer

import (
    "github.com/Shopify/sarama"
)

type Consumer struct {
    BrokersUrl []string
    Consumer   sarama.Consumer
}

func (c *Consumer) Close() error {
    return c.Consumer.Close()
}

func NewConsumer(brokersUrl []string) (con *Consumer, err error) {
    con = &Consumer{
        BrokersUrl: brokersUrl,
    }
    con.Consumer, err = sarama.NewConsumer(con.BrokersUrl, nil)

    return
}
