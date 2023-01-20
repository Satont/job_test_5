package amqp

import (
	qp "github.com/rabbitmq/amqp091-go"
)

func CreateAmqpConn(rabbitUrl string) (*qp.Connection, error) {
	conn, err := qp.Dial(rabbitUrl)

	return conn, err
}

func CreateChannel(conn *qp.Connection) (*qp.Channel, error) {
	ch, err := conn.Channel()
	return ch, err
}

func CreateQueue(ch *qp.Channel, queue string) (*qp.Queue, error) {
	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	return &q, err
}

func CreateConsume(ch *qp.Channel, queue string) (<-chan qp.Delivery, error) {
	c, err := ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	return c, err
}

func CreateMessage(body []byte) qp.Publishing {
	return qp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	}
}
