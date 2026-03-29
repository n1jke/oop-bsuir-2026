package application

import "fmt"

type WarehouseWorker interface {
	ProcessOrder()
	GetRest()
}

type HumanWarehouseWorker interface {
	WarehouseWorker
	AttendMeeting()
	SwingingTheLead()
}

type RobotWarehouseWorker interface {
	WarehouseWorker
	// may add specific methods
}

// ManageWarehouse - функция, которая работает со списком работников.
func ManageWarehouse(workers []WarehouseWorker) {
	fmt.Println("\n--- Warehouse Shift Started ---")

	for _, worker := range workers {
		worker.ProcessOrder()
		worker.GetRest()
	}
}
