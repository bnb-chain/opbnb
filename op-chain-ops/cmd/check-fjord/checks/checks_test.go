package checks

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/require"
)

func TestIsReceiptNotReady(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "ethereum not found", err: ethereum.NotFound, want: true},
		{name: "legacy not found message", err: errors.New("transaction not found"), want: true},
		{name: "transaction indexing", err: errors.New("transaction indexing is in progress"), want: true},
		{name: "other RPC error", err: errors.New("connection reset"), want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, isReceiptNotReady(test.err))
		})
	}
}
