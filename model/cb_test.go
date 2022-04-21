package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCBet(t *testing.T) {
	type args struct {
		playerID int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"OK",
			args{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				got := NewCBet(tt.args.playerID)
				assert.Equal(t, tt.args.playerID, got.PlayerId)
			})
		})
	}
}
