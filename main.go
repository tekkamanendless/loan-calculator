package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	loan := Loan{
		Amount:    39125,
		Rate:      4.99,
		Months:    20 * 12,
		Payment:   187.48,
		StartDate: "2019-10-01",
	}
	extras := []Extra{
		{
			Frequency: "monthly",
			Amount:    52.52,
			StartDate: "2019-10-01",
			EndDate:   "2099-01-01",
		},
		{
			Frequency: "once",
			StartDate: "2020-11-10",
			Amount:    12000,
		},
		{
			Frequency: "once",
			StartDate: "2020-12-20",
			Amount:    5000,
		},
	}
	schedule := Calculate(loan, extras)
	fmt.Printf("%3s   %10s   %8s  %8s   %8s\n", "", "date", "principal", "interest", "remaining")
	for i, payment := range schedule {
		fmt.Printf("%3d   %10s   %8.2f   %8.2f   %8.2f\n", i, payment.Date.Format("2006-01-02"), payment.Principal, payment.Interest, payment.Remaining)
	}
	fmt.Printf("----\n")
	fmt.Printf("Total interest paid: %8.2f\n", schedule[len(schedule)-1].InterestPaid)
}

type Loan struct {
	Amount    float64
	Rate      float64
	Months    int
	Payment   float64
	StartDate string
}

type Extra struct {
	Frequency string
	StartDate string
	EndDate   string
	Amount    float64
}

type Payment struct {
	Date          time.Time
	Principal     float64
	Interest      float64
	Remaining     float64
	PrincipalPaid float64
	InterestPaid  float64
}

func Calculate(loan Loan, extras []Extra) []Payment {
	startDate, err := time.Parse("2006-01-02", loan.StartDate)
	if err != nil {
		panic(err)
	}
	rate := loan.Rate / 100.0

	var payments []Payment
	amountRemaining := loan.Amount
	var interestPaid float64
	var principalPaid float64
	for m := 0; m < loan.Months; m++ {
		currentDate := startDate.AddDate(0, m, 0)

		currentDateString := currentDate.Format("2006-01-02")
		for e, extra := range extras {
			if extra.Frequency == "once" && strings.Compare(extra.StartDate, currentDateString) <= 0 {
				extraDate, err := time.Parse("2006-01-02", extra.StartDate)
				if err != nil {
					panic(err)
				}

				amountRemaining -= extra.Amount
				principalPaid += extra.Amount
				payment := Payment{
					Date:          extraDate,
					Principal:     extra.Amount,
					Interest:      0,
					Remaining:     amountRemaining,
					PrincipalPaid: principalPaid,
					InterestPaid:  interestPaid,
				}
				payments = append(payments, payment)

				extras = append(extras[:e], extras[e+1:]...)
				break
			}
		}

		if amountRemaining <= 0 {
			break
		}

		var extraPrincipal float64
		for _, extra := range extras {
			if extra.Frequency == "monthly" && strings.Compare(currentDateString, extra.StartDate) >= 0 && strings.Compare(currentDateString, extra.EndDate) < 0 {
				extraPrincipal = extra.Amount
			}
		}

		interest := amountRemaining * (rate / 12.0)
		principal := loan.Payment - interest
		if principal < 0 {
			principal = 0
		}
		principal += extraPrincipal
		if amountRemaining < principal {
			principal = amountRemaining
		}
		amountRemaining -= principal
		interestPaid += interest
		principalPaid += principal

		payment := Payment{
			Date:          currentDate,
			Interest:      interest,
			Principal:     principal,
			Remaining:     amountRemaining,
			PrincipalPaid: principalPaid,
			InterestPaid:  interestPaid,
		}
		payments = append(payments, payment)
	}
	return payments
}
