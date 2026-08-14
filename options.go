package main

type Focusable interface {
	Focus()
	Blur()
}

// func clamp(value, min, max int) int {
// 	if value < min {
// 		return min
// 	}
// 	if value > max {
// 		return max
// 	}
// 	return value
// }
