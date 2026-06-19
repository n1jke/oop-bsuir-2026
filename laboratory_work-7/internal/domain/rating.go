package domain

type Rating struct {
	points      uint
	reviewCount uint
}

func (r *Rating) AddReview(points uint) {
	r.points += points
	r.reviewCount++
}

func (r *Rating) GetRating() float64 {
	if r.reviewCount == 0 {
		return 0
	}

	return float64(r.points) / float64(r.reviewCount)
}
