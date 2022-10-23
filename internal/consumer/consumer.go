package consumer

import (
	"context"
	"log"

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

type ConsumerHandler struct {
	Ready chan bool
}

func (c *ConsumerGroup) Close() error {
	return c.Consumer.Close()
}

func (c *Consumer) Close() error {
	return c.Consumer.Close()
}

func NewConsumerGroup(brokersUrl []string, groupId string) (con *ConsumerGroup, err error) {

	version, err := sarama.ParseKafkaVersion("")

	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = version
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

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

func (c *ConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	return c.Consumer.Consume(ctx, topics, handler)
}

func IsLastMessage(cons sarama.PartitionConsumer, msg *sarama.ConsumerMessage) bool {
	return cons.HighWaterMarkOffset() == msg.Offset+1
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
	for {
		select {
		case message := <-claim.Messages():
			log.Printf(
				"Message claimed: value = %s, timestamp = %v, topic = %s",
				string(message.Value),
				message.Timestamp,
				message.Topic,
			)
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
