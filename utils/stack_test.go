package utils

import (
	"testing"
)

type STestStruct struct {
	eval bool
	fun  func() bool
}

var (
	spushTests = []STestStruct{
		{true,
			func() bool {
				exArr := []int{1, 2, 3}
				var s Stack[int]
				s.Push(1)
				s.Push(2)
				s.Push(3)
				return IsSame(exArr, s.GetArray())
			},
		},
		{false,
			func() bool {
				exArr := []int{1, 2}
				var s Stack[int]
				s.Push(1)
				s.Push(2)
				s.Push(3)
				s.Push(4)
				s.Push(5)
				return IsSame(exArr, s.GetArray())
			},
		},
	}

	spopTests = []STestStruct{
		{false,
			func() bool {
				exArr := []int{4}
				var s Stack[int]
				s.Push(1)
				s.Push(2)
				s.Push(3)
				s.Push(4)
				s.Pop()
				s.Pop()
				s.Pop()
				return IsSame(exArr, s.GetArray())
			},
		},
		{true,
			func() bool {
				exArr := []int{1, 2}
				var s Stack[int]
				s.Push(1)
				s.Push(2)
				s.Push(3)
				s.Push(4)
				s.Push(5)
				s.Pop()
				s.Pop()
				s.Pop()
				return IsSame(exArr, s.GetArray())
			},
		},
	}
)

func TestSPush(t *testing.T) {
	for _, v := range spushTests {
		if v.eval != v.fun() {
			t.Errorf("Receieved %v, expected : %v for %v", v.eval, v.fun(), v)
		}
	}
}

func TestSPop(t *testing.T) {
	for _, v := range spopTests {
		if v.eval != v.fun() {
			t.Errorf("Receieved %v, expected : %v for %v", v.fun(), v.eval, v)
		}
	}
}
