package utils

import (
	"fmt"

	"github.com/jaevor/go-nanoid"
)

var standardNanoID, _ = nanoid.Custom("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 18)

func NewLeaseID() string {
	return fmt.Sprintf("lea_%s", standardNanoID())
}
