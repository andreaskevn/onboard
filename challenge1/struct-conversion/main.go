package main

import "fmt"

type ConversionResult struct {
	Currency string
	Amount   float64
}

func main() {
	var rupiah float64
	var res []ConversionResult
	change := map[string]float64{
		"USD": 15000,
		"EUR": 16000,
		"JPY": 140,
		"SGD": 11000,
	}

	fmt.Print("Masukkan mata uang dalam Rupiah: ")
	fmt.Scanln(&rupiah)

	if rupiah == 0 {
		fmt.Println("Amount dalam rupiah tidak boleh 0!")
	} else {
		fmt.Printf("Konversi dari %.0f IDR: \n", rupiah)
		for currency, val := range change {
			multiply := rupiah / val
			exchange := ConversionResult{
				Currency: currency,
				Amount:   multiply,
			}

			res = append(res, exchange)
		}

		for _, val := range res {
			fmt.Printf("%s : %.3f\n", val.Currency, val.Amount)
		}
	}
}
