package application

import "github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"

type StockLoader interface {
	LoadCargoInfo() ([]domain.CargoInfo, error)
	LoadTransportInfo() ([]domain.TransportInfo, error)
}

type OrderRequestSource interface {
	RequestOrder(cargo []domain.CargoInfo, transport []domain.TransportInfo) (*ClientResponse, error)
}
