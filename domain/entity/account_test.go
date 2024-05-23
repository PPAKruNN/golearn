package entity

import (
	"crypto/sha256"
	"testing"
	"time"
)

func mockAccount(balance int) *Account {
	secret := sha256.New()

	return NewAccount(-1, balance, "Mock", "15799999970", secret, time.Now().UTC())
}

func TestIsValid(t *testing.T) {

	const (
		positiveAmount int = 100
		negativeAmount int = -100
	)

	t.Run("Should create account with positive amount", func(t *testing.T) {
		account := mockAccount(positiveAmount)
		isValid, err := account.IsValid()

		if err != nil || !isValid {
			t.Errorf("Account should be created with positive amount")
		}
	})

	t.Run("Should NOT create account with negative amount", func(t *testing.T) {
		account := mockAccount(negativeAmount)
		isValid, err := account.IsValid()

		if err == nil || isValid {
			t.Errorf("Account should NOT be created with negative amount")
		}
	})
}

func TestTransferTo(t *testing.T) {

	// acc1 := mockAccount(100)
	// acc2 := mockAccount(200)
	//
	// transfer, err := acc1.TransferTo(acc2, 100)

}
