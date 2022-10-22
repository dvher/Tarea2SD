package consumer

import (
	"context"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	BrokersUrl []string
	Consumer   sarama.Consumer
}

type ConsumerGroup struct {
	BrokersUrl []string
	Consumer   sarama.ConsumerGroup
}

func (c *ConsumerGroup) Close() error {
	return c.Consumer.Close()
}

func (c *Consumer) Close() error {
	return c.Consumer.Close()
}

func NewConsumerGroup(brokersUrl []string, groupId string) (con *ConsumerGroup, err error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(brokersUrl, groupId, config)

	if err != nil {
		return nil, err
	}

	con = &ConsumerGroup{
		BrokersUrl: brokersUrl,
		Consumer:   consumer,
	}

	return con, nil
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

func (c *ConsumerGroup) Consume(topic string, handler sarama.ConsumerGroupHandler) error {
	return c.Consumer.Consume(context.Background(), []string{topic}, handler)
}
