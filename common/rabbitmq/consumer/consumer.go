package consumer

import (
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	done          chan struct{}
	addr          string
	exchange      string
	queue         string
	routingKey    string
	consumerTag   string
	autoDelete    bool
	handler       func([]byte) error
	delivery      <-chan amqp.Delivery
}

// NewConsumer 创建消费者
func NewConsumer(addr, exchange, queue, routingKey string, handler func([]byte) error) *Consumer {
	hostname, _ := os.Hostname()
	return &Consumer{
		addr:        addr,
		exchange:    exchange,
		queue:       queue,
		routingKey:  routingKey,
		consumerTag: hostname,
		autoDelete:  false,
		handler:     handler,
		done:        make(chan struct{}),
	}
}

func (c *Consumer) Start() error {
	if err := c.Run(); err != nil {
		return err
	}
	go c.ReConnect()
	return nil
}

func (c *Consumer) Stop() {
	close(c.done)

	if !c.conn.IsClosed() {
		// 关闭 SubMsg message delivery
		if err := c.channel.Cancel(c.consumerTag, true); err != nil {
			log.Println("rabbitmq consumer - channel cancel failed: ", err)
		}

		if err := c.conn.Close(); err != nil {
			log.Println("rabbitmq consumer - connection close failed: ", err)
		}
	}
}

func (c *Consumer) Run() (err error) {
	if c.conn, err = amqp.Dial(c.addr); err != nil {
		return err
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		c.conn.Close()
		return err
	}

	defer func() {
		if err != nil {
			c.channel.Close()
			c.conn.Close()
		}
	}()

	// 声明一个主要使用的 exchange
	err = c.channel.ExchangeDeclare(
		c.exchange,
		"x-delayed-message",
		true,
		c.autoDelete,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "fanout",
		})
	if err != nil {
		return err
	}

	// 声明一个延时队列, 延时消息就是要发送到这里
	q, err := c.channel.QueueDeclare(c.queue, false, c.autoDelete, false, false, nil)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(q.Name, "", c.exchange, false, nil)
	if err != nil {
		return err
	}

	c.delivery, err = c.channel.Consume(
		q.Name, c.consumerTag, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go c.Handle()

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error))
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))
	return
}

func (c *Consumer) ReConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				log.Println("rabbitmq consumer - connection NotifyClose: ", err)
			}
		case err := <-c.channelNotify:
			if err != nil {
				log.Println("rabbitmq consumer - channel NotifyClose: ", err)
			}
		case <-c.done:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			// close message delivery
			if err := c.channel.Cancel(c.consumerTag, true); err != nil {
				log.Println("rabbitmq consumer - channel cancel failed: ", err)
			}

			if err := c.conn.Close(); err != nil {
				log.Println("rabbitmq consumer - channel cancel failed: ", err)
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for err := range c.channelNotify {
			log.Println(err)
		}
		for err := range c.connNotify {
			log.Println(err)
		}

	quit:
		for {
			select {
			case <-c.done:
				return
			default:
				log.Println("rabbitmq consumer - reconnect")

				if err := c.Run(); err != nil {
					log.Println("rabbitmq consumer - failCheck: ", err)
					// sleep 15s reconnect
					time.Sleep(time.Second * 15)
					continue
				}
				break quit
			}
		}
	}
}

func (c *Consumer) Handle() {
	for d := range c.delivery {
		go func(delivery amqp.Delivery) {
			if err := c.handler(delivery.Body); err != nil {
				// 重新入队，否则未确认的消息会持续占用内存，这里的操作取决于你的实现，你可以当出错之后并直接丢弃也是可以的
				_ = delivery.Reject(true)
			} else {
				_ = delivery.Ack(false)
			}
		}(d)
	}
}
