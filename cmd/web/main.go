// +build js

package main

import (
	"fmt"
	"syscall/js"

	loancalc "github.com/tekkamanendless/loan-calculator"
)

func main() {
	fmt.Printf("Go Web Assembly\n")

	validateLoanFunction := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := map[string]interface{}{
			"success": true,
		}
		if len(args) != 1 {
			result["success"] = false
			result["message"] = "Invalid number of arguments."
			return result
		}
		input := args[0].String()
		_, err := loancalc.ParseLoan(input)
		if err != nil {
			result["success"] = false
			result["message"] = err.Error()
			return result
		}
		return result
	})

	validateExtraFunction := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := map[string]interface{}{
			"success": true,
		}
		if len(args) != 1 {
			result["success"] = false
			result["message"] = "Invalid number of arguments."
			return result
		}
		input := args[0].String()
		_, err := loancalc.ParseExtra(input)
		if err != nil {
			result["success"] = false
			result["message"] = err.Error()
			return result
		}
		return result
	})

	calculateFunction := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := map[string]interface{}{
			"success": true,
		}
		if len(args) == 0 {
			result["success"] = false
			result["message"] = "Invalid number of arguments."
			return result
		}
		loanString := args[0].String()
		var extraStrings []string
		args = args[1:]
		for _, arg := range args {
			extraString := arg.String()
			extraStrings = append(extraStrings, extraString)
		}
		fmt.Printf("loanString: %s\n", loanString)
		fmt.Printf("extraStrings: %v\n", extraStrings)

		var loan *loancalc.Loan
		var extras []loancalc.Extra

		fmt.Printf("Parsing loan...\n")
		{
			v, err := loancalc.ParseLoan(loanString)
			if err != nil {
				result["success"] = false
				result["message"] = err.Error()
				return result
			}
			loan = v
		}
		for _, extraString := range extraStrings {
			fmt.Printf("Parsing extra...\n")
			v, err := loancalc.ParseExtra(extraString)
			if err != nil {
				result["success"] = false
				result["message"] = err.Error()
				return result
			}
			extras = append(extras, *v)
		}

		fmt.Printf("Calculating...\n")
		schedule := loancalc.Calculate(*loan, extras)
		fmt.Printf("Success.\n")
		{
			list := []interface{}{}
			for _, payment := range schedule {
				item := map[string]interface{}{
					"date":          payment.Date.Format("2006-01-02"),
					"principal":     payment.Principal,
					"interest":      payment.Interest,
					"principalPaid": payment.PrincipalPaid,
					"interestPaid":  payment.InterestPaid,
					"remaining":     payment.Remaining,
				}
				list = append(list, item)
			}
			result["data"] = list
		}
		return result
	})

	js.Global().Set("loancalc_calculate", calculateFunction)
	js.Global().Set("loancalc_validateLoan", validateLoanFunction)
	js.Global().Set("loancalc_validateExtra", validateExtraFunction)
	<-make(chan bool)
}
