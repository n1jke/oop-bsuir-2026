package domain

import "errors"

var ErrTransportNotFound = errors.New("transport not found")

type TransportFactory interface {
	Create(name TransportType) (Transport, error)
}

type CatalogTransportFactory struct {
	catalog []TransportInfo
}

func NewCatalogTransportFactory(catalog []TransportInfo) *CatalogTransportFactory {
	copied := make([]TransportInfo, len(catalog))
	copy(copied, catalog)

	return &CatalogTransportFactory{catalog: copied}
}

func (f *CatalogTransportFactory) Create(name TransportType) (Transport, error) {
	for i := range f.catalog {
		info := f.catalog[i]
		if info.Name() != name {
			continue
		}

		return newTransportFromInfo(info), nil
	}

	return nil, ErrTransportNotFound
}

type transportModel struct {
	TransportInfo
}

func newTransportFromInfo(info TransportInfo) Transport {
	return &transportModel{
		TransportInfo: info,
	}
}
