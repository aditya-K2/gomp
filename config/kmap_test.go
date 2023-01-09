package config

import (
	"testing"
)

func TestGetAsciiValue(t *testing.T) {
	for k, v := range SPECIAL_KEYS {
		result, err := GetAsciiValue(k)
		if result != v {
			t.Errorf("Values From KMAP Failed")
		} else if err != nil {
			t.Errorf("Error Received: %v", err)
		}
	}
	for i := 65; i < (65 + 26); i++ {
		result, err := GetAsciiValue(string(rune(i)))
		if err != nil {
			t.Errorf("Error Received! %v", err)
		} else if result != i {
			t.Errorf("Invalid Value Received for small alphabets Result Received: %d expecting %d", result, i)
		}
	}
	for i := 97; i < (97 + 26); i++ {
		result, err := GetAsciiValue(string(rune(i)))
		if err != nil {
			t.Errorf("Error Received: %v", err)
		} else if result != i {
			t.Errorf("Invalid Value Received for Big alphabets Result Received: %d expecting %d", result, i)
		}
	}
	result, err := GetAsciiValue("")
	if err == nil || result != -1 {
		t.Errorf("Wrong Result Received for Control Characters %v %v", result, err)
	}
}
