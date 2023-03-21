package producer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Producer struct {
	Addr       string
	Exchange   string
	RoutingKey string
	Queue      string

	conn       *amqp.Connection
	Channel    *amqp.Channel
	done       chan bool
	connErr    chan error
	channelErr chan *amqp.Error
}

func NewProducer(addr string, exchange string, queue string, key string) *Producer {
	return &Producer{
		Addr:       addr,
		Exchange:   exchange,
		RoutingKey: key,
		Queue:      queue,
		done:       make(chan bool),
		connErr:    make(chan error),
		channelErr: make(chan *amqp.Error),
	}
}

// connect 连接到mq服务器
func (c *Producer) Connect() error {
	var err error
	if c.conn, err = amqp.Dial(c.Addr); err != nil {
		return err
	}

	if c.Channel, err = c.conn.Channel(); err != nil {
		_ = c.Close()
		return err
	}

	// 声明一个主要使用的 exchange
	err = c.Channel.ExchangeDeclare(
		c.Exchange, "x-delayed-message", true, false, false, false, amqp.Table{
			"x-delayed-type": "fanout",
		})
	if err != nil {
		return err
	}
	// 声明一个延时队列, 延时消息就是要发送到这里
	q, err := c.Channel.QueueDeclare(c.Queue, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = c.Channel.QueueBind(q.Name, "", c.Exchange, false, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Producer) Close() error {
	close(c.done)

	if !c.conn.IsClosed() {
		if err := c.conn.Close(); err != nil {
			log.Print("rabbitmq producer - connection close failed: ", err)
			return err
		}
	}
	return nil
}

// publish 发送消息至mq
func (c *Producer) Publish(body []byte, delay int64) error {
	publising := amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	}

	if delay >= 0 {
		publising.Headers = amqp.Table{
			"x-delay": delay,
		}
	}

	err := c.Channel.Publish(c.Exchange, c.RoutingKey, false, false, publising)
	if err != nil {
		switch v := err.(type) {
		case *amqp.Error:
			c.channelErr <- v
		default:
			c.connErr <- v
		}
	}
	return err
}

// PublishJSON 将对象JSON格式化后发送消息
func (c *Producer) PublishJSON(body interface{}, delay int64) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return c.Publish(data, delay)
}

// watchconn 监控mq的连接状态
func (c *Producer) WatchConnect() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case err := <-c.connErr:
			log.Printf("rabbitmq producer - connection notify close: %s", err.Error())
			c.ReConnect()
		case err := <-c.channelErr:
			log.Printf("rabbitmq producer - channel notify close: %s", err.Error())
			c.ReConnect()
		case <-ticker.C:
			c.ReConnect()

		case <-c.done:
			log.Print("auto detect connection is done")
			return

		}
	}

}

// ReConnect 根据当前链接状态判断是否需要重新连接，如果连接异常则尝试重新连接
func (c *Producer) ReConnect() {
	if c.conn == nil || (c.conn != nil && c.conn.IsClosed()) {
		log.Printf("rabbitmq connection is closed try to reconnect")
		if err := c.Connect(); err != nil {
			log.Printf("rabbitmq reconnect failed: %s", err.Error())
		} else {
			log.Printf("rabbitmq reconnect succeeded")
		}
	}
}
