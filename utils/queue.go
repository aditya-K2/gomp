package utils

type Queue[T any] struct {
	arr []T
}

func (s *Queue[T]) GetArray() []T {
	return s.arr
}

func (s *Queue[T]) Push(item T) {
	s.arr = append(s.arr, item)
}

func (s *Queue[T]) Pop() {
	s.arr = s.arr[1:]
}

func (s *Queue[T]) Top() T {
	return s.arr[0]
}

func (s *Queue[T]) Empty() bool {
	return len(s.arr) == 0
}
