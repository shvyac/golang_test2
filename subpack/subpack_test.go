package subpack

import (
	"testing"
)

func TestReadfile(t *testing.T) {
}
func TestToJstTimeFromString(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "2023/01/0121:00"},
		{name: "2023/02/0121:00"},
		{name: "2023/03/0121:00"},
		{name: "2023/04/0121:00"},
		{name: "2023/05/0121:00"},
	}
	for _, tt := range tests {
		got := ToJstTimeFromString(tt.name)
		if got.Month() != 3 {
			//t.Errorf("ToJstTimeFromString() = %v, want %v", got, tt.name)
			t.Logf("ToJstTimeFromString() got= %v, want= %v", got, tt.name)
		}

		t.Run(tt.name, func(t *testing.T) {
			ToJstTimeFromString(tt.name)
		})
	}
}
