package loancalc

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
			if i >= len(parts) {
				return nil, fmt.Errorf("missing amount")
			}
			part = strings.TrimSpace(parts[i])
			v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
			if err != nil {
				return nil, err
			}
			loan.Amount = v
		case "rate":
			i++
			if i >= len(parts) {
				return nil, fmt.Errorf("missing rate")
			}
			part = strings.TrimSpace(parts[i])
			v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
			if err != nil {
				return nil, err
			}
			loan.Rate = v
		case "months":
			i++
			if i >= len(parts) {
				return nil, fmt.Errorf("missing months")
			}
			part = strings.TrimSpace(parts[i])
			v, err := strconv.ParseInt(strings.ReplaceAll(part, ",", ""), 10, 64)
			if err != nil {
				return nil, err
			}
			loan.Months = int(v)
		case "payment":
			i++
			if i >= len(parts) {
				return nil, fmt.Errorf("missing payment")
			}
			part = strings.TrimSpace(parts[i])
			v, err := strconv.ParseFloat(strings.ReplaceAll(part, ",", ""), 64)
			if err != nil {
				return nil, err
			}
			loan.Payment = v
		case "starting":
			i++
			if i >= len(parts) {
				return nil, fmt.Errorf("missing starting")
			}
			part = strings.TrimSpace(parts[i])
			v, err := time.Parse(DateFormat, part)
			if err != nil {
				return nil, err
			}
			loan.StartDate = v
		}
	}
	if loan.StartDate.IsZero() {
		loan.StartDate = time.Now()
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
				if i >= len(parts) {
					return nil, fmt.Errorf("missing count")
				}
				part = strings.TrimSpace(parts[i])
				v, err := strconv.ParseInt(strings.ReplaceAll(part, ",", ""), 10, 64)
				if err != nil {
					return nil, err
				}
				if v < 0 {
					return nil, fmt.Errorf("count is less than zero")
				}
				extra.Count = int(v)
			case "on", "starting":
				i++
				if i >= len(parts) {
					return nil, fmt.Errorf("missing %s", part)
				}
				part = strings.TrimSpace(parts[i])
				v, err := time.Parse(DateFormat, part)
				if err != nil {
					return nil, err
				}
				extra.StartDate = v
			case "ending":
				i++
				if i >= len(parts) {
					return nil, fmt.Errorf("missing ending")
				}
				part = strings.TrimSpace(parts[i])
				v, err := time.Parse(DateFormat, part)
				if err != nil {
					return nil, err
				}
				extra.EndDate = v
			}
		}
	}

	if extra.Frequency == "" {
		return nil, fmt.Errorf("missing frequency")
	}
	if !extra.EndDate.IsZero() && extra.Count > 0 {
		return nil, fmt.Errorf("cannot have both end date and count")
	}
	if extra.Frequency == "once" && extra.StartDate.IsZero() {
		return nil, fmt.Errorf("missing start date")
	}

	return extra, nil
}
