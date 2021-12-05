package internal

import (
	"constraints"
	"fmt"
)

func GetOrCreate[K comparable, V any](dict map[K]V, key K, create func(K) V) V {
	value, exists := dict[key]
	if !exists {
		value = create(key)
		dict[key] = value
	}
	return value
}

func Append[T any](slice []T, values []T) []T {
	dst := make([]T, len(slice))
	copy(dst, slice)
	return append(dst, values...)
}

func RemoveIndex[T any](slice []T, index int) []T {
	return Append[T](slice[:index], slice[index+1:])
}

func Map[T, R any](list []T, transform func(T) R) []R {
	result := make([]R, len(list))
	for index, item := range list {
		result[index] = transform(item)
	}
	return result
}

func Filter[T any](list []T, cond func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range list {
		if cond(item) {
			result = append(result, item)
		}
	}
	return result
}

func Find[T any](list []T, cond func(T) bool) int {
	for index, item := range list {
		if cond(item) {
			return index
		}
	}
	return -1
}

func FindValue[T comparable](list []T, value T) int {
	for index, item := range list {
		if item == value {
			return index
		}
	}
	return -1
}

func Reduce[T any](values []T, init T, action func(acc T, value T) T) T {
	acc := init
	for _, value := range values {
		acc = action(acc, value)

	}
	return acc
}

func Min[T constraints.Ordered](init T, values ...T) T {
	min := func(lhs T, rhs T) T {
		if lhs < rhs {
			return lhs
		} else {
			return rhs
		}
	}

	return Reduce(values, init, min)
}

func Skip[T any](slice []T, n int) int {
	return Min(len(slice), n)
}

func Demo() {
	s := []string{"one", "two", "three", "four"}
	fmt.Println(s[Skip(s, 3):])

	var str interface{} = 123
	value, ok := str.(string)
	fmt.Println(ok)
	fmt.Println(value)
}
