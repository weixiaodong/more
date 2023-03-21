package rabbitmq

import (
	"sync"

	"github.com/weixiaodong/more/common/config"
	"github.com/weixiaodong/more/common/rabbitmq/consumer"
	"github.com/weixiaodong/more/common/rabbitmq/producer"
)

var (
	producerMaps sync.Map
)

func GetProducer(exchange string) *producer.Producer {

	if v, ok := producerMaps.Load(exchange); ok {
		return v.(*producer.Producer)
	}

	p := producer.NewProducer(
		config.GetRabbitMQURL(),
		"exchange",
		"queue",
		"",
	)
	err := p.Connect()
	if err != nil {
		panic(err)
	}
	go p.WatchConnect()

	producerMaps.Store(exchange, p)
	return p
}

func PublishHello(msg interface{}, delay int64) error {
	p := GetProducer("hello")
	return p.PublishJSON(msg, delay)
}

func StartRabbitMQConsumer(h func([]byte) error) {
	c := consumer.NewConsumer(config.GetRabbitMQURL(), "exchange", "queue", "", h)
	c.Start()
}
