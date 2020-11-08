package accounts

import (
	"errors"
	"fmt"
)

// Account struct.
type Account struct {
	owner   string
	balance int
}

var errNoMoney = errors.New("Cant Withdraw")

//NewAccount creates func
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

//Deposit x amount on your balance
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

//Balance show account balance
func (a Account) Balance() int {
	return a.balance
}

//Withdraw x amount from your account
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}

//ChangeOwner of the account
func (a *Account) ChangeOwner(newOwners string) {
	a.owner = newOwners
}

//Owner from the account
func (a *Account) Owner() string {
	return a.owner
}

func (a *Account) String() string {
	return fmt.Sprint(a.owner, "'s account.\nHas:", a.balance)
}
