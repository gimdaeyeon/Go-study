package main

import (
	"fmt"
	"learngo/accounts"
)

func main() {
	account := accounts.NewAccount("머연")
	account.Deposit(10)

	fmt.Println(account)

}
