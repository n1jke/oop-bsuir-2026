package application

import (
	"errors"
	"log/slog"

	"github.com/n1jke/oop-bsuir-2025/laboratory_work-5/internal/domain"
)

var (
	ErrNilStockRepository  = errors.New("stock repository must not be nil")
	ErrNilOrderSource      = errors.New("order request source must not be nil")
	ErrNilLogger           = errors.New("logger must not be nil")
	ErrTransportNotInStock = errors.New("selected transport is not present in stock")
	ErrCargoNotInStock     = errors.New("selected cargo is not present in stock")
)

type StockLoader interface {
	LoadCargoInfo() ([]domain.CargoInfo, error)
	LoadTransportInfo() ([]domain.TransportInfo, error)
}

type OrderRequestSource interface {
	RequestOrder(cargo []domain.CargoInfo, transport []domain.TransportInfo) (*ClientResponse, error)
}

type ClientResponse struct {
	Transport domain.Transport
	Dist      float64
	Content   []domain.ProductBatch
}

type ServiceResponse struct {
	Order *domain.Order
	Cost  float64
}

func (c ClientResponse) ToOrder() (*domain.Order, error) {
	return domain.NewOrder(c.Dist, c.Transport, c.Content)
}

type LogisticService struct {
	logger *slog.Logger
	stock  StockLoader
	client OrderRequestSource
}

type Option func(l *LogisticService) error

func NewLogisticService(opts ...Option) (*LogisticService, error) {
	srv := &LogisticService{
		logger: slog.Default(),
	}

	for _, opt := range opts {
		if err := opt(srv); err != nil {
			return srv, err
		}
	}

	if srv.stock == nil {
		return nil, ErrNilStockRepository
	}

	if srv.client == nil {
		return nil, ErrNilOrderSource
	}

	return srv, nil
}

func WithLogger(logger *slog.Logger) Option {
	return func(l *LogisticService) error {
		if logger == nil {
			return ErrNilLogger
		}

		l.logger = logger

		return nil
	}
}

func WithStock(stock StockLoader) Option {
	return func(l *LogisticService) error {
		if stock == nil {
			return ErrNilStockRepository
		}

		l.stock = stock

		return nil
	}
}

func WithClient(client OrderRequestSource) Option {
	return func(l *LogisticService) error {
		if client == nil {
			return ErrNilOrderSource
		}

		l.client = client

		return nil
	}
}

func (l *LogisticService) Process() (*ServiceResponse, error) {
	l.logger.Info("Processing order logistic")

	cargoConfig, err := l.stock.LoadCargoInfo()
	if err != nil {
		l.logger.Error("Error loading cargo info", "error", err)
		return nil, err
	}

	transportConfig, err := l.stock.LoadTransportInfo()
	if err != nil {
		l.logger.Error("Error loading transport info", "error", err)
		return nil, err
	}

	l.logger.Info("Get all configuration")

	resp, err := l.client.RequestOrder(cargoConfig, transportConfig)
	if err != nil {
		l.logger.Error("Error requesting order", "error", err)
		return nil, err
	}

	order, err := resp.ToOrder()
	if err != nil {
		l.logger.Error("Error parsing order response", "error", err)
		return nil, err
	}

	response := &ServiceResponse{
		Order: order,
		Cost:  order.CalculateCost(),
	}

	return response, nil
}
