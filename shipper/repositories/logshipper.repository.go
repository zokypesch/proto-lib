package repositories

type LogShipperRepository interface {
	Send(topic string, value []byte) error
}
