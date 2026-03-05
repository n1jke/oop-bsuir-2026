package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInactiveAccount      = errors.New("account is not active")
	ErrNegativeAmount       = errors.New("amount of money is negative")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrBonusBalanceTooLow   = errors.New("bonus balance is too low")
	ErrInvalidBonusRedeem   = errors.New("bonus points must be positive")
	ErrSavingsOnlyOperation = errors.New("operation is allowed only for savings account")
)

// Account - bank account entity.
type Account struct {
	id       uuid.UUID
	number   string
	clientID uuid.UUID
	balance  Money
	status   AccountStatus
}

type SavingsAccount struct {
	Account
	BonusProgram BonusProgram
}

type CreditAccount struct {
	Account
	overdraftLimit Money
}

// AccountStatus - value object for bank account status.
type AccountStatus int

const (
	Active AccountStatus = iota + 1
	Frozen
	Closed
)

func NewAccount(number string, clientID uuid.UUID, curr Currency) *Account {
	return &Account{
		id:       uuid.New(),
		number:   number,
		clientID: clientID,
		balance:  NewMoney(0, curr),
		status:   Active,
	}
}

func NewSavingsAccount(number string, accountUUID uuid.UUID, curr Currency, tier BonusTier) *SavingsAccount {
	return &SavingsAccount{
		Account:      *NewAccount(number, accountUUID, curr),
		BonusProgram: NewBonusProgram(tier),
	}
}

func NewCreditAccount(number string, accountUUID uuid.UUID, curr Currency, overdraftLimit Money) *CreditAccount {
	return &CreditAccount{
		Account:        *NewAccount(number, accountUUID, curr),
		overdraftLimit: overdraftLimit,
	}
}

func (a *Account) ID() uuid.UUID {
	return a.id
}

func (a *Account) Number() string {
	return a.number
}

func (a *Account) ClientID() uuid.UUID {
	return a.clientID
}

func (a *Account) Balance() Money {
	return a.balance
}

func (a *Account) Status() AccountStatus {
	return a.status
}

func (a *Account) Deposit(value Money) error {
	if a.status != Active {
		return ErrInactiveAccount
	}

	if value.Amount() < 0 {
		return ErrNegativeAmount
	}

	a.balance.Add(value)

	return nil
}

func (s *SavingsAccount) Deposit(value Money) error {
	err := s.Account.Deposit(value)
	if err != nil {
		return err
	}

	s.BonusProgram.Accrue(value)

	return nil
}

func (s *SavingsAccount) BonusPoints() int {
	return s.BonusProgram.Bonus()
}

func (s *SavingsAccount) RedeemBonus(points int) error {
	if points <= 0 {
		return ErrInvalidBonusRedeem
	}

	if !s.BonusProgram.ApplyBonus(points) {
		return ErrBonusBalanceTooLow
	}

	return nil
}

func (a *Account) Withdraw(value Money) error {
	if err := a.CanWithdraw(value); err != nil {
		return err
	}

	a.balance.Sub(value)

	return nil
}

func (a *Account) CanWithdraw(value Money) error {
	if a.status != Active {
		return ErrInactiveAccount
	}

	if a.balance.amount < value.amount {
		return ErrInsufficientFunds
	}

	return nil
}

func (c *CreditAccount) CanWithdraw(value Money) error {
	if c.status != Active {
		return ErrInactiveAccount
	}

	if c.balance.amount+c.overdraftLimit.amount < value.amount {
		return ErrInsufficientFunds
	}

	return nil
}

func (c *CreditAccount) Withdraw(value Money) error {
	if err := c.CanWithdraw(value); err != nil {
		return err
	}

	c.balance.Sub(value)

	return nil
}

func (a *Account) IsActive() bool {
	return a.status == Active
}

func (a *Account) ChangeStatus(status AccountStatus) {
	a.status = status
}
