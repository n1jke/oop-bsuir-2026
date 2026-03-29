package application

import "strings"

func compareByCost(a, b Quote) int {
	switch {
	case a.Cost < b.Cost:
		return -1
	case a.Cost > b.Cost:
		return 1
	default:
		return 0
	}
}

func compareByDuration(a, b Quote) int {
	if a.Duration == b.Duration {
		return 0
	}

	if a.Duration < b.Duration {
		return -1
	}

	return 1
}

func compareByTransportName(a, b Quote) int {
	na := strings.ToLower(string(a.Transport.Name()))
	nb := strings.ToLower(string(b.Transport.Name()))

	return strings.Compare(na, nb)
}
