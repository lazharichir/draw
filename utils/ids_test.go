package utils_test

import (
	"testing"

	"github.com/lazharichir/draw/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewLeaseID(t *testing.T) {
	// Generate a new lease ID.
	leaseID := utils.NewLeaseID()

	// Test that the lease ID has the correct prefix.
	assert.Equal(t, "lea_", leaseID[:4])

	// Test that the lease ID has the correct length.
	assert.GreaterOrEqual(t, 20, len(leaseID))
}
