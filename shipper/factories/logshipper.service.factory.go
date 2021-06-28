package factories

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/prakerja/logger/shipper/repositories"
	"gitlab.com/prakerja/logger/shipper/services"
)

var logShipperServiceInstance *services.V1LogShipperService

func MakeLogShipper(producer *kafka.Producer, topic string) (services.LogShipperService, error) {
	if logShipperServiceInstance == nil {
		services.NewV1LogShipperService(
			repositories.NewKafkaLogShipperRepository(producer),
			topic,
		)
	}

	return logShipperServiceInstance, nil
}
