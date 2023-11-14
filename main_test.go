package main

import (
	"testing"
)

func TestToUnixTimeFromString(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "2023/06/0121:00"},
		{name: "2023/07/0121:00"},
		{name: "2023/08/0121:00"},
		{name: "2023/09/0121:00"},
		{name: "2023/10/0121:00"},
	}
	for _, tt := range tests {
		got := ToUnixTimeFromString(tt.name)
		//t.Errorf("ToJstTimeFromString() = %v, want %v", got, tt.name)
		t.Logf("ToUnixTimeFromString() got= %v, want= %v", got, tt.name)
		// t.Run(tt.name, func(t *testing.T) {
		// 	ToUnixTimeFromString(tt.name)
		// })
	}
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
		// t.Run(tt.name, func(t *testing.T) {
		// 	ToJstTimeFromString(tt.name)
		// })
	}
}
