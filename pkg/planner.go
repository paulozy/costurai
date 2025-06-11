package pkg

type PlannerPrice struct {
	MonthlyStandard int32
	MonthlyPro      int32
	YearlyStandard  int32
	YearlyPro       int32
}

func NewPlannerPrice() *PlannerPrice {
	return &PlannerPrice{
		MonthlyStandard: 349,
		MonthlyPro:      999,
		YearlyStandard:  3799,
		YearlyPro:       999,
	}
}
