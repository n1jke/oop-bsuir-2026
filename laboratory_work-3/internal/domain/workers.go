package domain

import "fmt"

// HumanManager - Человек.
type HumanManager struct{}

func (h HumanManager) ProcessOrder() {
	fmt.Println("Manager is processing logic...")
}

func (h HumanManager) AttendMeeting() {
	fmt.Println("Manager is boring at the meeting...")
}

func (h HumanManager) GetRest() {
	fmt.Println("Manager is taking a break...")
}

func (h HumanManager) SwingingTheLead() {
	fmt.Println("Manager is watching reels...")
}

// RobotPacker - Робот.
type RobotPacker struct {
	Model string
}

func (r RobotPacker) ProcessOrder() {
	fmt.Println("Robot " + r.Model + " is packing boxes...")
}

func (r RobotPacker) GetRest() {
	fmt.Println("Robot was taken for maintenance")
}
