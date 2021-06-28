package services_test

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/prakerja/logger/shipper/entities"
	"gitlab.com/prakerja/logger/shipper/repositories"
	"gitlab.com/prakerja/logger/shipper/services"
	"testing"
)

func TestShipLogUsingChannel(t *testing.T) {
	type args struct {
		request chan entities.ShipLogRequest
	}

	var testData = []struct {
		name string
		args args
	}{
		{
			name: "Test ship log using channel",
			args: args{
				request: func() chan entities.ShipLogRequest {
					ch := make(chan entities.ShipLogRequest)

					go func() {
						ch <- entities.ShipLogRequest{
							Action:   "http://localhost/test1",
							Todo:     "POST",
							Property: `{"Request":{"headers":"content-type: application/json\r\nx-internal-token: 1nIint3n4Lt0KenNy4\r\nuserID: 951706\r\nUser-Agent: PostmanRuntime/7.28.0\r\nAccept: /\r\nPostman-Token: ed090688-a6e5-4f7f-a386-00e19edaab7c\r\nHost: localhost:8080\r\nAccept-Encoding: gzip, deflate, br\r\nConnection: keep-alive\r\n\r\n","payload":"","others":{}},"Response":{"status_code":500,"content":"{\"errorCode\":\"ERR9999\",\"success\":false,\"message\":\"record not found\"}\n","others":{}}}`,
						}

						ch <- entities.ShipLogRequest{
							Action:   "http://localhost/test2",
							Todo:     "POST",
							Property: `{"Request":{"headers":"content-type: application/json\r\nx-internal-token: 1nIint3n4Lt0KenNy4\r\nuserID: 951706\r\nUser-Agent: PostmanRuntime/7.28.0\r\nAccept: /\r\nPostman-Token: ed090688-a6e5-4f7f-a386-00e19edaab7c\r\nHost: localhost:8080\r\nAccept-Encoding: gzip, deflate, br\r\nConnection: keep-alive\r\n\r\n","payload":"","others":{}},"Response":{"status_code":500,"content":"{\"errorCode\":\"ERR9999\",\"success\":false,\"message\":\"record not found\"}\n","others":{}}}`,
						}

						ch <- entities.ShipLogRequest{
							Action:   "http://localhost/test3",
							Todo:     "POST",
							Property: `{"Request":{"headers":"content-type: application/json\r\nx-internal-token: 1nIint3n4Lt0KenNy4\r\nuserID: 951706\r\nUser-Agent: PostmanRuntime/7.28.0\r\nAccept: /\r\nPostman-Token: ed090688-a6e5-4f7f-a386-00e19edaab7c\r\nHost: localhost:8080\r\nAccept-Encoding: gzip, deflate, br\r\nConnection: keep-alive\r\n\r\n","payload":"","others":{}},"Response":{"status_code":500,"content":"{\"errorCode\":\"ERR9999\",\"success\":false,\"message\":\"record not found\"}\n","others":{}}}`,
						}
						close(ch)
					}()

					return ch
				}(),
			},
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			producer, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": "172.28.15.189:9092",
			})
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			service := services.NewV1LogShipperService(
				repositories.NewKafkaLogShipperRepository(producer),
				"logging",
			)

			for e := range service.Send(tt.args.request) {
				if e != nil {
					t.Errorf("%v", e)
				}
			}
		})
	}
}


func TestShipSingleLog(t *testing.T) {
	type args struct {
		request entities.ShipLogRequest
	}

	var testData = []struct {
		name string
		args args
	}{
		{
			name: "Test ship single log",
			args: args{
				request: entities.ShipLogRequest{
					Action:   "http://localhost/test1",
					Todo:     "POST",
					Property: `{"Request":{"headers":"content-type: application/json\r\nx-internal-token: 1nIint3n4Lt0KenNy4\r\nuserID: 951706\r\nUser-Agent: PostmanRuntime/7.28.0\r\nAccept: /\r\nPostman-Token: ed090688-a6e5-4f7f-a386-00e19edaab7c\r\nHost: localhost:8080\r\nAccept-Encoding: gzip, deflate, br\r\nConnection: keep-alive\r\n\r\n","payload":"","others":{}},"Response":{"status_code":500,"content":"{\"errorCode\":\"ERR9999\",\"success\":false,\"message\":\"record not found\"}\n","others":{}}}`,
				},
			},
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			producer, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": "172.28.15.189:9092",
				"linger.ms": 5,
			})
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			service := services.NewV1LogShipperService(
				repositories.NewKafkaLogShipperRepository(producer),
				"logging",
			)

			service.SendSingle(tt.args.request)
		})
	}
}


func TestConsumer(t *testing.T) {
	type args struct {
	}

	var testData = []struct {
		name string
		args args
	}{
		{
			name: "Test consumer",
			args: args{
			},
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			c, err := kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": "172.28.15.189:9092",
				"group.id":          "myGroup1",
				"auto.offset.reset": "earliest",
			})

			if err != nil {
				panic(err)
			}

			c.SubscribeTopics([]string{"logging"}, nil)

			for {
				msg, err := c.ReadMessage(-1)
				if err == nil {
					fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				} else {
					// The client will automatically try to recover from all errors.
					fmt.Printf("Consumer error: %v (%v)\n", err, msg)
				}
			}

			c.Close()
		})
	}
}
