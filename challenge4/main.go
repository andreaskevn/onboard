package main

import (
	"context"
	"fmt"
	// "math"
	"math/rand"
	// "runtime"
	"sync"
	"time"
	// "github.com/creasty/defaults"
)

type Account struct {
	AccountNumber string
	Balance       int
}

var Accounts = map[string]*Account{
	"C001": {AccountNumber: "C001", Balance: 300000},
	"M002": {AccountNumber: "M002", Balance: 500000},
	"B003": {AccountNumber: "B003", Balance: 400000},
	"BC04": {AccountNumber: "BC04", Balance: 800000},
}

var mu sync.Mutex

type TransferRequest struct {
	ID          int
	FromAccount string
	ToAccount   string
	Amount      int
}

func TransferAsync(req TransferRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	resultChan := make(chan string)
	errChan := make(chan error)

	go func() {
		// random delay 1–5 detik
		delay := time.Duration(rand.Intn(7)+1) * time.Second
		fmt.Printf("TRX %d processing delay: %v\n", req.ID, delay)
		time.Sleep(delay)

		fromAcc := Accounts[req.FromAccount]
		toAcc := Accounts[req.ToAccount]

		if fromAcc == nil || toAcc == nil {
			errChan <- fmt.Errorf("TRX %d account not found", req.ID)
			return
		}

		// LOCK
		mu.Lock()
		defer mu.Unlock()

		// cek saldo
		if fromAcc.Balance < req.Amount {
			resultChan <- fmt.Sprintf(
				"[FAILED] TRX %d | %s -> %s | Amount: %d | Balance: %d",
				req.ID, req.FromAccount, req.ToAccount, req.Amount, fromAcc.Balance,
			)
			return
		}

		// update saldo
		fromAcc.Balance -= req.Amount
		toAcc.Balance += req.Amount

		resultChan <- fmt.Sprintf(
			"[SUCCESS] TRX %d | %s(%d) -> %s(%d)",
			req.ID,
			req.FromAccount, fromAcc.Balance,
			req.ToAccount, toAcc.Balance,
		)
		// defer mu.Unlock()
		// for msg := range resultChan {
		// 	fmt.Printf("hasil: %s", msg)
		// }
	}()

	select {
	case res := <-resultChan:
		fmt.Println(res)
	case err := <-errChan:
		fmt.Println("[ERROR]", err)
	case <-ctx.Done():
		fmt.Println("[TIMEOUT] TRX", req.ID)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	transactions := []TransferRequest{
		{ID: 1, FromAccount: "C001", ToAccount: "M002", Amount: 300000},
		{ID: 2, FromAccount: "C001", ToAccount: "B003", Amount: 600000},
		{ID: 3, FromAccount: "M002", ToAccount: "BC04", Amount: 200000},
		{ID: 4, FromAccount: "B003", ToAccount: "C001", Amount: 500000},
		{ID: 5, FromAccount: "BC04", ToAccount: "M002", Amount: 700000},
		{ID: 6, FromAccount: "C001", ToAccount: "M002", Amount: 400000},
		{ID: 7, FromAccount: "B003", ToAccount: "C001", Amount: 300000},
	}

	var wg sync.WaitGroup
	// wx := sync.WaitGroup/

	for _, trx := range transactions {
		wg.Add(1)
		// go TransferAsync(trx, &wg)
		go TransferAsync(trx, &wg)
	}

	wg.Wait()

	fmt.Println("\n=== FINAL BALANCE ===")
	totalBalance := []int{}

	for _, acc := range Accounts {
		fmt.Printf("%s: %d\n", acc.AccountNumber, acc.Balance)

		totalBalance = append(totalBalance, acc.Balance)
	}

	// hitung total
	sum := 0
	for _, val := range totalBalance {
		sum += val
	}

	fmt.Printf("Total Balance: %d\n", sum)
}
