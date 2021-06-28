package services

import (
	"encoding/json"
	"gitlab.com/prakerja/logger/shipper/entities"
	"gitlab.com/prakerja/logger/shipper/repositories"
)

type LogShipperService interface {
	Send(request <-chan entities.ShipLogRequest) <-chan error
	SendBatch(requests []entities.ShipLogRequest) <-chan error
	SendSingle(request entities.ShipLogRequest) <-chan error
}

type V1LogShipperService struct {
	repo                repositories.LogShipperRepository
	topic               string
}

func NewV1LogShipperService(repo repositories.LogShipperRepository, topic string) *V1LogShipperService {
	return &V1LogShipperService{repo: repo, topic: topic}
}

func (v V1LogShipperService) Send(requestCh <-chan entities.ShipLogRequest) <-chan error {
	ch := make(chan error)

	go func() {
		for r := range requestCh {
			bytes, err := json.Marshal(r)
			if err != nil {
				ch <- err
			} else {
				ch <- v.repo.Send(v.topic, bytes)
			}
		}

		close(ch)
	}()

	return ch
}

func (v V1LogShipperService) SendSingle(request entities.ShipLogRequest) <-chan error {
	ch := make(chan entities.ShipLogRequest, 1)
	ch <- request
	close(ch)

	return v.Send(ch)
}

func (v V1LogShipperService) SendBatch(requests []entities.ShipLogRequest) <-chan error {
	return v.Send(v.makeRequestChanFromBatch(requests))
}

func (v V1LogShipperService) makeRequestChanFromBatch(requests []entities.ShipLogRequest) <-chan entities.ShipLogRequest {
	ch := make(chan entities.ShipLogRequest)

	go func() {
		for _, r := range requests {
			ch <- r
		}
		close(ch)
	}()

	return ch
}
