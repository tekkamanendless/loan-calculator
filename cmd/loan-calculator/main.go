package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	loancalc "github.com/tekkamanendless/loan-calculator"
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
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Missing loan information.\n")
		help()
	}
	loan, err := loancalc.ParseLoan(args[0])
	if err != nil {
		fmt.Printf("Could not parse loan: %v", err)
		help()
	}
	if loan.StartDate == "" {
		loan.StartDate = time.Now().Format("2006-01-02")
	}

	var extras []loancalc.Extra
	for _, input := range args[1:] {
		extra, err := loancalc.ParseExtra(input)
		if err != nil {
			fmt.Printf("Could not parse extra: %v", err)
			help()
		}

		extras = append(extras, *extra)
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

	schedule := loancalc.Calculate(*loan, extras)

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
