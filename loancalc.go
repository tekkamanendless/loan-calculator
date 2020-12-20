package loancalc

import (
	"time"
)

// DateFormat is the format used for dates.
const DateFormat = "2006-01-02"

// Loan is a loan.
type Loan struct {
	Amount    float64
	Rate      float64
	Months    int
	Payment   float64
	StartDate time.Time
}

// Extra payment.
type Extra struct {
	Frequency string
	Count     int
	StartDate time.Time
	EndDate   time.Time
	Amount    float64
}

// Payment on the schedule.
type Payment struct {
	Date          time.Time
	Principal     float64
	Interest      float64
	Remaining     float64
	PrincipalPaid float64
	InterestPaid  float64
}

// Calculate a payment schedule.
func Calculate(loan Loan, extras []Extra) []Payment {
	for e := range extras {
		if extras[e].Count > 0 {
			extras[e].EndDate = loan.StartDate.AddDate(0, extras[e].Count, 0)
		}
	}

	rate := loan.Rate / 100.0

	var payments []Payment
	amountRemaining := loan.Amount
	var interestPaid float64
	var principalPaid float64
	for m := 0; m < loan.Months; m++ {
		currentDate := loan.StartDate.AddDate(0, m, 0)

		for e, extra := range extras {
			if extra.Frequency == "once" && !extra.StartDate.After(currentDate) {
				amountRemaining -= extra.Amount
				principalPaid += extra.Amount
				payment := Payment{
					Date:          extra.StartDate,
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
			if extra.Frequency == "monthly" && (extra.StartDate.IsZero() || !currentDate.Before(extra.StartDate)) && (extra.EndDate.IsZero() || currentDate.Before(extra.EndDate)) {
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
		if m+1 == loan.Months {
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
