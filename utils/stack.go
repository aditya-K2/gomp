package utils

type Stack[T any] struct {
	arr []T
}

func (s *Stack[T]) Push(item T) {
	s.arr = append(s.arr, item)
}

func (s *Stack[T]) Pop() {
	s.arr = s.arr[0 : len(s.arr)-1]
}

func (s *Stack[T]) Top() T {
	return s.arr[len(s.arr)-1]
}

func (s *Stack[T]) Empty() bool {
	return len(s.arr) == 0
}

func (s Stack[T]) GetArray() []T {
	return s.arr
}
