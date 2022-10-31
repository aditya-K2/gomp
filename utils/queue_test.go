package utils

import (
	"testing"
)

type QTestStruct struct {
	eval bool
	fun  func() bool
}

var (
	pushTests = []QTestStruct{
		{true,
			func() bool {
				exArr := []int{1, 2, 3}
				var q Queue[int]
				q.Push(1)
				q.Push(2)
				q.Push(3)
				return IsSame(exArr, q.GetArray())
			},
		},
		{false,
			func() bool {
				exArr := []int{1, 2}
				var q Queue[int]
				q.Push(1)
				q.Push(2)
				q.Push(3)
				q.Push(4)
				q.Push(5)
				return IsSame(exArr, q.GetArray())
			},
		},
	}

	popTests = []QTestStruct{
		{true,
			func() bool {
				exArr := []int{4}
				var q Queue[int]
				q.Push(1)
				q.Push(2)
				q.Push(3)
				q.Push(4)
				q.Pop()
				q.Pop()
				q.Pop()
				return IsSame(exArr, q.GetArray())
			},
		},
		{false,
			func() bool {
				exArr := []int{1, 2}
				var q Queue[int]
				q.Push(1)
				q.Push(2)
				q.Push(3)
				q.Push(4)
				q.Push(5)
				q.Pop()
				q.Pop()
				q.Pop()
				return IsSame(exArr, q.GetArray())
			},
		},
	}
)

func TestQPush(t *testing.T) {
	for _, v := range pushTests {
		if v.eval != v.fun() {
			t.Errorf("Receieved %v, expected : %v for %v", v.eval, v.fun(), v)
		}
	}
}

func TestQPop(t *testing.T) {
	for _, v := range popTests {
		if v.eval != v.fun() {
			t.Errorf("Receieved %v, expected : %v for %v", v.fun(), v.eval, v)
		}
	}
}
