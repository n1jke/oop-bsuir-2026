package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/application"
	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/application/services"
	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/domain"
	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/infrastructure"
)

func main() {
	bank := domain.NewBank(uuid.New(), "BelInvest", "234567")
	branch := domain.NewBranch(uuid.New(), bank.ID(), "Minsk, Independence avenue")

	client := domain.NewClient(uuid.New(), "BM9817236", "Ivan Pumpalumpa")
	employee := domain.NewEmployee(uuid.New(), branch.ID(), "manager")
	loan := domain.NewLoan(uuid.New(), client.ID(), 5000)
	loan.Approve()

	srcAccount := domain.NewAccount(uuid.New(), "BY009872", client.ID(), "BYN")
	dstAccount := domain.NewAccount(uuid.New(), "BY006789", client.ID(), "BYN")
	card := domain.NewCard(uuid.New(), "9387656789", srcAccount.ID())
	atm := domain.NewATM(uuid.New(), branch.ID(), 10000)

	accountStore := infrastructure.NewMemoryAccountStorage()
	eventStore := infrastructure.NewMemoryEventStorage()
	accountRepo := services.NewAccountRepository(accountStore)
	eventService := services.NewEventService(eventStore)
	paymentService := services.NewPaymentService(accountRepo)
	useCase := application.NewTransferUseCase(paymentService, eventService)

	_ = accountRepo.Create(srcAccount)
	_ = accountRepo.Create(dstAccount)

	_ = paymentService.Deposit(2000, srcAccount.ID())

	tx := domain.NewTransaction(uuid.New(), srcAccount.ID(), dstAccount.ID(), 838, "BYN")
	if err := useCase.Execute(tx); err != nil {
		log.Printf("transaction failed: %v", err)
		return
	}

	if card.IsActive() {
		card.Block()
		card.Activate()
	}

	_ = atm.Withdraw(300)
	atm.Deposit(100)

	events := eventService.QueryAll()
	fmt.Printf("bank=%s branch_opened=%t employee=%s employed=%t loan_approved=%t tx_completed=%t src_balance=%s dst_balance=%s events=%d\n",
		bank.Name(),
		branch.IsOpened(),
		employee.Position(),
		employee.IsEmployed(),
		loan.IsApproved(),
		tx.IsCompleted(),
		srcAccount.Balance(),
		dstAccount.Balance(),
		len(events),
	)

	fmt.Println("Events")

	for _, e := range events {
		fmt.Println(e)
	}
}
