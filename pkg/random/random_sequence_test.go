package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandSequence(t *testing.T) {
	type args struct {
		seed int
	}

	tests := []struct {
		name   string
		args   args
		length int
	}{
		{name: "success-string-10", args: args{seed: 10}, length: 10},
		{name: "success-string-20", args: args{seed: 20}, length: 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandSequence(tt.args.seed)
			got2 := RandSequence(tt.args.seed)
			got3 := RandSequence(tt.args.seed)

			assert.Equal(t, tt.length, len(got))
			assert.Equal(t, tt.length, len(got2))
			assert.Equal(t, tt.length, len(got3))
			assert.NotEqual(t, got, got2)
			assert.NotEqual(t, got, got3)
			assert.NotEqual(t, got2, got3)
		})
	}
}
