package domain

import (
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

func NewDatePolicy(createdAt, updatedAT, expiresAT time.Time) *DatePolicy {
	return &DatePolicy{
		createdAt: createdAt,
		updatedAT: updatedAT,
		expiresAT: expiresAT,
	}
}

func (d *DatePolicy) CreatedAt() time.Time { return d.createdAt }

func (d *DatePolicy) UpdatedAt() time.Time { return d.updatedAT }

func (d *DatePolicy) ExpiresAt() time.Time { return d.expiresAT }

func (d *DatePolicy) ChangeExpireDate(expDate time.Time) error {
	now := time.Now()
	if now.After(expDate) {
		return ErrExpireDateInPast
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

	if fromID == toID {
		return nil, ErrSelfExchange
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

func CreateExchangeRequest(id, ownedBookID, fromID, toID uuid.UUID, status ExchangeStatus, dateInfo *DatePolicy, note string,
) (*ExchangeRequest, error) {
	return &ExchangeRequest{
		id:          id,
		ownedBookID: ownedBookID,
		fromID:      fromID,
		toID:        toID,
		status:      status,
		dateInfo:    dateInfo,
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

func (er *ExchangeRequest) Accept() error {
	if er.status != Pending {
		return ErrInvalidExchangeTransition
	}

	er.status = Accepted
	er.dateInfo.Update()

	return nil
}

func (er *ExchangeRequest) Reject() error {
	if er.status != Pending {
		return ErrInvalidExchangeTransition
	}

	er.status = Rejected
	er.dateInfo.Update()

	return nil
}

func (er *ExchangeRequest) Complete() error {
	if er.status != Accepted {
		return ErrInvalidExchangeTransition
	}

	er.status = Completed
	er.dateInfo.Update()

	return nil
}

func (er *ExchangeRequest) Cancel() error {
	if er.status == Completed || er.status == Canceled {
		return ErrInvalidExchangeTransition
	}

	er.status = Canceled
	er.dateInfo.Update()

	return nil
}

func (er *ExchangeRequest) Note() string {
	return er.note
}

func (er *ExchangeRequest) DateInfo() DatePolicy {
	return *er.dateInfo
}

func (er *ExchangeRequest) ChangeExpireDate(expDate time.Time) error {
	return er.dateInfo.ChangeExpireDate(expDate)
}
