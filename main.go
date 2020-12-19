package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func help() {
	fmt.Printf(strings.TrimSpace(`
Usage:
loan-calculator <loan> [<payment> [...]]

<loan> is of the form:
   amount <amount> rate <rate> months <months> payment <payment> [starting <date>]

<payment> is of the form:
   <amount> monthly [starting <date>] [ending <date>|count <count>]
	<amount> once on <date>

Example:
   loan-calculator 'amount 39,125.00 rate 4.99 months 240 payment 187.48 starting 2020-10-01' '52.52 monthly' '12,000 once on 2020-11-10' '5,000 once on 2020-12-20'
`) + "\n")
	os.Exit(1)
}

func main() {
	var loan Loan
	var extras []Extra

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Missing loan information.\n")
		help()
	}
	{
		input := strings.TrimSpace(args[0])
		parts := strings.Split(input, " ")
		for i := 0; i < len(parts); i++ {
			part := strings.TrimSpace(parts[i])
			switch part {
			case "amount":
				i++
				part = strings.TrimSpace(parts[i])
				v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
				if err != nil {
					help()
					panic(err)
				}
				loan.Amount = v
			case "rate":
				i++
				part = strings.TrimSpace(parts[i])
				v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
				if err != nil {
					help()
					panic(err)
				}
				loan.Rate = v
			case "months":
				i++
				part = strings.TrimSpace(parts[i])
				v, err := strconv.ParseInt(strings.ReplaceAll(part, ",", ""), 10, 64)
				if err != nil {
					help()
					panic(err)
				}
				loan.Months = int(v)
			case "payment":
				i++
				part = strings.TrimSpace(parts[i])
				v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
				if err != nil {
					help()
					panic(err)
				}
				loan.Payment = v
			case "starting":
				i++
				part = strings.TrimSpace(parts[i])
				loan.StartDate = part
			}
		}
	}
	if loan.StartDate == "" {
		loan.StartDate = time.Now().Format("2006-01-02")
	}

	for _, input := range args[1:] {
		extra := Extra{}

		input := strings.TrimSpace(input)
		parts := strings.Split(input, " ")
		for i := 0; i < len(parts); i++ {
			part := strings.TrimSpace(parts[i])
			if i == 0 {
				v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
				if err != nil {
					help()
					panic(err)
				}
				extra.Amount = v
			} else if i == 1 {
				switch part {
				case "once", "monthly":
					extra.Frequency = part
				default:
					help()
				}
			} else {
				switch part {
				case "count":
					i++
					part = strings.TrimSpace(parts[i])
					v, err := strconv.ParseInt(strings.ReplaceAll(part, ",", ""), 10, 64)
					if err != nil {
						help()
						panic(err)
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

		extras = append(extras, extra)
	}
	for e := range extras {
		if extras[e].StartDate == "" {
			extras[e].StartDate = loan.StartDate
		}
		if extras[e].Count > 0 {
			if extras[e].EndDate != "" {
				help()
			}
			date, err := time.Parse("2006-01-02", extras[e].StartDate)
			if err != nil {
				panic(err)
			}
			date = date.AddDate(0, extras[e].Count, 0)
			extras[e].EndDate = date.Format("2006-01-02")
		}
	}

	p := message.NewPrinter(language.English)

	fmt.Printf("Loan:\n")
	fmt.Printf("   Amount:   %12s $\n", p.Sprintf("%0.2f", loan.Amount))
	fmt.Printf("   Rate:     %12s %%\n", p.Sprintf("%0.2f", loan.Rate))
	fmt.Printf("   Months:   %12s\n", p.Sprintf("%0d", loan.Months))
	fmt.Printf("   Payment:  %12s $\n", p.Sprintf("%0.2f", loan.Payment))
	fmt.Printf("   Starting: %12s\n", loan.StartDate)
	fmt.Printf("\n")

	for _, extra := range extras {
		fmt.Printf("Extra payment:\n")
		fmt.Printf("   Amount:   %12s $ %s\n", p.Sprintf("%0.2f", extra.Amount), extra.Frequency)
		if extra.StartDate != "" {
			fmt.Printf("   Starting: %12s\n", extra.StartDate)
		}
		if extra.EndDate != "" {
			fmt.Printf("   Ending:   %12s\n", extra.EndDate)
		}
		fmt.Printf("\n")
	}

	schedule := Calculate(loan, extras)

	fmt.Printf("Schedule:\n")
	fmt.Printf("%3s   %10s   %10s   %10s   %12s\n", "", "date", "principal", "interest", "remaining")
	for i, payment := range schedule {
		fmt.Printf("%3d   %10s   %10s   %10s   %12s\n", i, payment.Date.Format("2006-01-02"), p.Sprintf("%0.2f", payment.Principal), p.Sprintf("%0.2f", payment.Interest), p.Sprintf("%0.2f", payment.Remaining))
	}
	if len(schedule) > 0 {
		fmt.Printf("----\n")
		fmt.Printf("Total interest paid: %12s\n", p.Sprintf("%0.2f", schedule[len(schedule)-1].InterestPaid))
	}
}

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
