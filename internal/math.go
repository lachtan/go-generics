package internal

import "constraints"

func Max[T constraints.Ordered](first T, others ...T) T {
	max := first
	for _, value := range others {
		if value > max {
			max = value
		}
	}
	return max
}
