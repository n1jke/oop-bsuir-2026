package main

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/application"
	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/application/services"
	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/domain"
	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/infrastructure"
)

func main() {
	// setup entities
	bank := domain.NewBank(uuid.New(), "BelInvest", "234567")
	branch := domain.NewBranch(uuid.New(), bank.ID(), "Minsk, Independence avenue")

	client := domain.NewClient(uuid.New(), "BM9817236", "Ivan Pumpalumpa")
	employee := domain.NewEmployee(branch.ID(), "manager")
	loan := domain.NewLoan(uuid.New(), client.ID(), 5000)
	loan.Approve()

	srcAccount := domain.NewAccount("BY009872", client.ID(), domain.Currency("BYN"))
	savingsAccount := domain.NewSavingsAccount("BY006789", client.ID(), domain.Currency("BYN"), domain.Gold)
	creditAccount := domain.NewCreditAccount("BY007654", client.ID(), domain.Currency("BYN"), domain.NewMoney(1500, domain.Currency("BYN")))
	card := domain.NewCard(uuid.New(), "9387656789", srcAccount.ID())
	atm := domain.NewATM(uuid.New(), branch.ID(), 10000)

	// infra
	accountStore := infrastructure.NewMemoryAccountStorage()
	eventStore := infrastructure.NewMemoryEventStorage()
	paymentService := services.NewPaymentService(accountStore)

	// usecasess
	transferUseCase := application.NewTransferUseCase(paymentService, eventStore)
	bonusUseCase := application.NewBonusRedeemUseCase(accountStore, eventStore)

	// storage setup
	_ = accountStore.Save(srcAccount.ID(), srcAccount)
	_ = accountStore.Save(savingsAccount.ID(), savingsAccount)
	_ = accountStore.Save(creditAccount.ID(), creditAccount)

	_ = paymentService.Deposit(domain.NewMoney(2000, domain.Currency("BYN")), srcAccount.ID())
	_ = paymentService.Deposit(domain.NewMoney(500, domain.Currency("BYN")), savingsAccount.ID())

	// transactions
	successTx := domain.NewTransaction(
		uuid.New(),
		srcAccount.ID(),
		savingsAccount.ID(),
		domain.NewMoney(838, domain.Currency("BYN")),
	)
	if err := transferUseCase.Execute(successTx); err != nil {
		fmt.Printf("success transfer UNexpectedly failed: %v\n", err)
	}

	creditTx := domain.NewTransaction(
		uuid.New(),
		creditAccount.ID(),
		srcAccount.ID(),
		domain.NewMoney(700, domain.Currency("BYN")),
	)
	if err := transferUseCase.Execute(creditTx); err != nil {
		fmt.Printf("credit transfer UNexpectedly failed: %v\n", err)
	}

	failedTx := domain.NewTransaction(
		uuid.New(),
		srcAccount.ID(),
		savingsAccount.ID(),
		domain.NewMoney(999999, domain.Currency("BYN")),
	)
	if err := transferUseCase.Execute(failedTx); err != nil {
		fmt.Printf("Expected failed transfer: %v\n", err)
	}

	fmt.Println()

	// bonuses
	if err := bonusUseCase.Execute(application.BonusRedeemCommand{
		AccountID: savingsAccount.ID(),
		Points:    20,
	}); err != nil {
		fmt.Printf("bonus redeem failed: %v\n", err)
	}

	if err := bonusUseCase.Execute(application.BonusRedeemCommand{
		AccountID: srcAccount.ID(),
		Points:    10,
	}); err != nil {
		fmt.Printf("Expected non-savings bonus failure: %v\n", err)
	}

	fmt.Println()

	// card & atm
	if card.IsActive() {
		card.Block()
		card.Activate()
	}

	_ = atm.Withdraw(300)
	atm.Deposit(100)

	// events query
	events := eventStore.QueryAll()
	fmt.Printf("bank=%s branch_opened=%t employee=%s employed=%t loan_approved=%t tx_completed=%t src_balance=%s savings_balance=%s credit_balance=%s savings_bonus=%d events=%d\n",
		bank.Name(),
		branch.IsOpened(),
		employee.Position(),
		employee.IsEmployed(),
		loan.IsApproved(),
		successTx.IsCompleted(),
		srcAccount.Balance(),
		savingsAccount.Balance(),
		creditAccount.Balance(),
		savingsAccount.BonusPoints(),
		len(events),
	)

	fmt.Println()
	fmt.Println("Events")

	for _, e := range events {
		fmt.Println(e)
	}
}
