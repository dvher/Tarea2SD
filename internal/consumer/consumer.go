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

    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true

    con = &Consumer{
        BrokersUrl: brokersUrl,
    }
    con.Consumer, err = sarama.NewConsumer(con.BrokersUrl, config)

    return
}

func (c *Consumer) Consume(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
    return c.Consumer.ConsumePartition(topic, partition, offset)
}

func (c *Consumer) ConsumeSinceLast(topic string, partition int32) (sarama.PartitionConsumer, error) {
    return c.Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
}

func (c *Consumer) ConsumeFromBeginning(topic string, partition int32) (sarama.PartitionConsumer, error) {
    return c.Consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
}
