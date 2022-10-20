package producer

import (
    "github.com/Shopify/sarama"
)

type Producer struct {
    BrokerList []string
    Producer   sarama.SyncProducer
}

func (p *Producer) Close() error {
    return p.Producer.Close()
}

func NewProducer(brokerList []string) (prod *Producer, err error) {
    prod = &Producer{
        BrokerList: brokerList,
    }

    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true
    config.Producer.Partitioner = sarama.NewManualPartitioner

    prod.Producer, err = sarama.NewSyncProducer(brokerList, config)
    if err != nil {
        return nil, err
    }

    return
}

func (p *Producer) SendMessage(topic string, part int32, message []byte) (partition int32, offset int64, err error) {
    msg := &sarama.ProducerMessage{
        Topic:     topic,
        Partition: part,
        Value:     sarama.StringEncoder(message),
    }

    partition, offset, err = p.Producer.SendMessage(msg)
    return
}
