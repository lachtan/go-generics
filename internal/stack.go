package internal

type Stack[T any] struct {
	vals []T
}

func (stack *Stack[T]) Push(val T) {
	stack.vals = append(stack.vals, val)
}

func (stack *Stack[T]) Pop() (T, bool) {
	if len(stack.vals) == 0 {
		var zero T
		return zero, false
	}
	top := stack.vals[len(stack.vals)-1]
	stack.vals = stack.vals[:len(stack.vals)-1]
	return top, true
}

type XStack[T any] []T

func (stack *XStack[T]) Push(val T) {
	*stack = append(*stack, val)
}

func (stack *XStack[T]) Pop() (T, bool) {
	if len(*stack) == 0 {
		var zero T
		return zero, false
	}
	high := len(*stack) - 1
	top := (*stack)[high]
	*stack = (*stack)[:high]
	return top, true
}

func Push[T any](stack *[]T, val T) {
	*stack = append(*stack, val)
}

func Pop[T any](stack *[]T) (T, bool) {
	if len(*stack) == 0 {
		var zero T
		return zero, false
	} else {
		high := len(*stack) - 1
		top := (*stack)[high]
		*stack = (*stack)[:high]
		return top, true
	}
}
