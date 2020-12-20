package loancalc

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Loan is a loan.
type Loan struct {
	Amount    float64
	Rate      float64
	Months    int
	Payment   float64
	StartDate string
}

// Extra payment.
type Extra struct {
	Frequency string
	Count     int
	StartDate string
	EndDate   string
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

// ParseLoan parses a string and returns a Loan.
//
// Format:
//    amount <amount> rate <rate> months <months> payment <payment> [starting <date>]
func ParseLoan(input string) (*Loan, error) {
	loan := &Loan{}
	input = strings.TrimSpace(input)
	parts := strings.Split(input, " ")
	for i := 0; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		switch part {
		case "amount":
			i++
			part = strings.TrimSpace(parts[i])
			v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
			if err != nil {
				return nil, err
			}
			loan.Amount = v
		case "rate":
			i++
			part = strings.TrimSpace(parts[i])
			v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
			if err != nil {
				return nil, err
			}
			loan.Rate = v
		case "months":
			i++
			part = strings.TrimSpace(parts[i])
			v, err := strconv.ParseInt(strings.ReplaceAll(part, ",", ""), 10, 64)
			if err != nil {
				return nil, err
			}
			loan.Months = int(v)
		case "payment":
			i++
			part = strings.TrimSpace(parts[i])
			v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
			if err != nil {
				return nil, err
			}
			loan.Payment = v
		case "starting":
			i++
			part = strings.TrimSpace(parts[i])
			loan.StartDate = part
		}
	}
	return loan, nil
}

// ParseExtra parses a string and returns an Extra.
//
// Format:
//    <amount> monthly [starting <date>] [ending <date>|count <count>]
//    <amount> once on <date>
func ParseExtra(input string) (*Extra, error) {
	extra := &Extra{}

	input = strings.TrimSpace(input)
	parts := strings.Split(input, " ")
	for i := 0; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		if i == 0 {
			v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
			if err != nil {
				return nil, err
			}
			extra.Amount = v
		} else if i == 1 {
			switch part {
			case "once", "monthly":
				extra.Frequency = part
			default:
				return nil, fmt.Errorf("invalid frequency: %s", part)
			}
		} else {
			switch part {
			case "count":
				i++
				part = strings.TrimSpace(parts[i])
				v, err := strconv.ParseInt(strings.ReplaceAll(part, ",", ""), 10, 64)
				if err != nil {
					return nil, err
				}
				extra.Count = int(v)
			case "on":
				i++
				part = strings.TrimSpace(parts[i])
				extra.StartDate = part
			case "starting":
				i++
				part = strings.TrimSpace(parts[i])
				extra.StartDate = part
			case "ending":
				i++
				part = strings.TrimSpace(parts[i])
				extra.EndDate = part
			}
		}
	}

	return extra, nil
}

// Calculate a payment schedule.
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
			if extra.Frequency == "monthly" && (extra.StartDate == "" || strings.Compare(currentDateString, extra.StartDate) >= 0) && (extra.EndDate == "" || strings.Compare(currentDateString, extra.EndDate) < 0) {
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
