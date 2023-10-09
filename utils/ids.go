package utils

import (
	"fmt"

	"github.com/jaevor/go-nanoid"
)

var eightNanoID, _ = nanoid.Custom("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 12)

// var twelveNanoID, _ = nanoid.Custom("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 12)
var eighteenNanoID, _ = nanoid.Custom("ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz123456789", 18)

func NewLeaseID() string {
	return fmt.Sprintf("lea_%s", eighteenNanoID())
}

func NewUserID() string {
	return fmt.Sprintf("usr_%s", eightNanoID())
}

func NewVerificationTokenSignup() string {
	return eighteenNanoID()
}

func NewVerificationTokenSignin() string {
	return eighteenNanoID()
}

func NewVerificationTokenChangeEmail() string {
	return eighteenNanoID()
}

func NewVerificationTokenResetPassword() string {
	return eighteenNanoID()
}
