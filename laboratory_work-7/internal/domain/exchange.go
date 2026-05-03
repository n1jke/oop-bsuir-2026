package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ExchangeStatus int

const (
	Pending ExchangeStatus = iota
	Accepted
	Rejected
	Completed
	Canceled
)

type DatePolicy struct {
	createdAt time.Time
	updatedAT time.Time
	expiresAT time.Time
}

func (d *DatePolicy) ChangeExpireDate(expDate time.Time) error {
	now := time.Now()
	if now.After(expDate) {
		return errors.New("oekf") // todo move to custom error or const err with fmt.Errorf
	}

	d.updatedAT = now
	d.expiresAT = expDate
	return nil
}

func (d *DatePolicy) Update() {
	d.updatedAT = time.Now()
}

type ExchangeRequest struct {
	id          uuid.UUID
	ownedBookID uuid.UUID
	fromID      uuid.UUID
	toID        uuid.UUID
	status      ExchangeStatus
	dateInfo    *DatePolicy
	note        string
}

func NewExchangeRequest(ownedBookID, fromID, toID uuid.UUID, d *DatePolicy, note string) (*ExchangeRequest, error) {
	if d == nil {
		return nil, ErrDatePolicyNotConfigure
	}
	return &ExchangeRequest{
		id:          uuid.New(),
		ownedBookID: ownedBookID,
		fromID:      fromID,
		toID:        toID,
		status:      Pending,
		dateInfo:    d,
		note:        note,
	}, nil
}

func (er *ExchangeRequest) ID() uuid.UUID {
	return er.id
}

func (er *ExchangeRequest) OwnedBookID() uuid.UUID {
	return er.ownedBookID
}

func (er *ExchangeRequest) FromID() uuid.UUID {
	return er.fromID
}

func (er *ExchangeRequest) ToID() uuid.UUID {
	return er.toID
}

func (er *ExchangeRequest) Status() ExchangeStatus {
	return er.status
}

// fine tune & improve state machine
func (e *ExchangeRequest) Accept() error {
	if e.status != Pending {
		return errors.New("") // create custom error using snipper goerr
	}
	e.status = Accepted
	e.dateInfo.Update()
	return nil
}

func (e *ExchangeRequest) Reject() error {
	if e.status != Pending {
		return errors.New("") // create custom error using snipper goerr
	}
	e.status = Rejected
	e.dateInfo.Update()
	return nil
}

func (er *ExchangeRequest) DateInfo() DatePolicy {
	// todo: maybe return copy?
	return *er.dateInfo
}

func (er *ExchangeRequest) ChangeExpireDate(expDate time.Time) error {
	return er.dateInfo.ChangeExpireDate(expDate)
}
