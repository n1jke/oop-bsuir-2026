package domain

type BonusProgram struct {
	points int
	tier   BonusTier
}

type BonusTier int

const (
	Bronze BonusTier = iota + 1
	Silver
	Gold
)

func NewBonusProgram(tier BonusTier) BonusProgram {
	return BonusProgram{tier: tier}
}

func (b *BonusProgram) Accrue(value Money) {
	b.points += value.amount * int(b.tier) / 100
}

func (b *BonusProgram) ApplyBonus(points int) bool {
	if b.points < points {
		return false
	}

	b.points -= points

	return true
}

func (b BonusProgram) Bonus() int {
	return b.points
}
