// package main

// import (
// 	"fmt"
// 	"strings"
// 	// "strconv"
// 	// "strings"
// )

// type Card struct {
// 	Type   string
// 	Prefix []string
// 	Length []int
// }

// func main() {
// 	var cardNum string

// 	var cardRules = []Card{
// 		{
// 			Type:   "China UnionPay",
// 			Prefix: []string{"62"},
// 			Length: []int{16, 17, 18, 19},
// 		},
// 		{
// 			Type:   "Switch",
// 			Prefix: []string{"4903", "4905", "4911", "4936", "564182", "633110", "6333", "6759"},
// 			Length: []int{16, 18, 19},
// 		},
// 	}

// }

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Card struct {
	Type   string
	Prefix []string
	Length []int
}

func main() {
	var nomor int
	fmt.Println("Menu Cek Kartu")
	fmt.Println("1. Bulk Input")
	fmt.Println("2. Single Input")

	fmt.Print("Masukkan nomor menu: ")
	fmt.Scanln(&nomor)

	var cardRules = []Card{
		{
			Type:   "China UnionPay",
			Prefix: []string{"62"},
			Length: []int{16, 17, 18, 19},
		},
		{
			Type:   "Switch",
			Prefix: []string{"4903", "4905", "4911", "4936", "564182", "633110", "6333", "6759"},
			Length: []int{16, 18, 19},
		},
	}

	switch nomor {
	case 1:
		fmt.Println("Masukkan nomor kartu:")

		scanner := bufio.NewScanner(os.Stdin)
		var inputList []string

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				break
			}
			inputList = append(inputList, line)
		}

		fmt.Println("\nHasil Pengecekan")

		for _, cardNum := range inputList {
			found := false

			for _, rule := range cardRules {
				for _, pref := range rule.Prefix {
					if strings.HasPrefix(cardNum, pref) {
						for _, length := range rule.Length {
							if len(cardNum) == length {
								fmt.Printf("Nomor %s adalah %s\n", cardNum, rule.Type)
								found = true
								break
							}
						}
					}
					if found {
						break
					}
				}
				if found {
					break
				}
			}

			if !found {
				fmt.Printf("Nomor %s Jenis kartu tidak dikenali\n", cardNum)
			}
		}
		// break

	case 2:
		var cardNum string
		fmt.Print("Masukkan nomor kartu: ")
		fmt.Scanln(&cardNum)

		find := false

		for _, val := range cardRules {
			for _, pref := range val.Prefix {
				if strings.HasPrefix(cardNum, pref) {
					for _, length := range val.Length {
						if len(cardNum) == length {
							fmt.Printf("Nomor %s adalah %s dengan panjang %d", cardNum, val.Type, len(cardNum))
							find = true
							break
						}
					}

					if !find {
						fmt.Println("Jenis kartu tidak dikenali")
						break
					}
				}
			}
		}
		break
	}
}
