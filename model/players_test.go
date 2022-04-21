package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayer(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"OK",
			args{id: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				got := NewPlayer(tt.args.id)
				assert.Equal(t, tt.args.id, got.ID)
			})
		})
	}
}
