package number

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestByteCountDecimal(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		want string
	}{
		{"bytes", 42, "42 B"},
		{"kilobytes", 1024, "1.0 kB"},
		{"megabytes", 1048576, "1.0 MB"},
		{"gigabytes", 1073741824, "1.1 GB"},
		{"terabytes", 1099511627776, "1.1 TB"},
		{"petabytes", 1125899906842624, "1.1 PB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ByteCountDecimal(tt.n)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestWithComma(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		want string
	}{
		{"zero", 0, "0"},
		{"positive number", 123456789, "123,456,789"},
		{"negative number", -987654321, "-987,654,321"},
		{"small number", 42, "42"},
		{"large number", 1234567890123456789, "1,234,567,890,123,456,789"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithComma(tt.n)
			require.Equal(t, tt.want, got)
		})
	}
}
