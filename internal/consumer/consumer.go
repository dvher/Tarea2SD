package consumer

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	brokersUrl []string
	consumer   sarama.Consumer
}

type ConsumerGroup struct {
	brokersUrl []string
	consumer   sarama.ConsumerGroup
}

type ConsumerHandler struct {
	Ready chan bool
}

func (c *ConsumerGroup) Close() error {
	return c.consumer.Close()
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func NewConsumerGroup(brokersUrl []string, groupId string, initialOffset int64) (con *ConsumerGroup, err error) {

	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = initialOffset

	consumer, err := sarama.NewConsumerGroup(brokersUrl, groupId, config)

	if err != nil {
		return nil, err
	}

	con = &ConsumerGroup{
		brokersUrl: brokersUrl,
		consumer:   consumer,
	}

	return con, nil
}

func NewConsumer(brokersUrl []string) (con *Consumer, err error) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	con = &Consumer{
		brokersUrl: brokersUrl,
	}
	con.consumer, err = sarama.NewConsumer(con.brokersUrl, config)

	return
}

func (c *Consumer) Consume(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	return c.consumer.ConsumePartition(topic, partition, offset)
}

func (c *Consumer) ConsumeSinceLast(topic string, partition int32) (sarama.PartitionConsumer, error) {
	return c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
}

func (c *Consumer) ConsumeFromBeginning(topic string, partition int32) (sarama.PartitionConsumer, error) {
	return c.consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
}

func (c *Consumer) Partitions(topic string) ([]int32, error) {
	return c.consumer.Partitions(topic)
}

func (c *ConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	return c.consumer.Consume(ctx, topics, handler)
}

func (consumer *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.Ready)
	return nil
}

func (consumer *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29

LOOP:
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, part = %d", string(message.Value), message.Timestamp, message.Topic, message.Partition)
			session.MarkMessage(message, "")
			log.Println(claim.HighWaterMarkOffset(), " ", message.Offset)

			if claim.HighWaterMarkOffset() == message.Offset+1 {
				break LOOP
			}

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}

	return nil
}

func IsLastMessage(cons sarama.PartitionConsumer, msg *sarama.ConsumerMessage) bool {
	return cons.HighWaterMarkOffset() == msg.Offset+1
}
