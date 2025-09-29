/*
Copyright © 2025 lixw
*/
package ternary

func IFElse[T any](ok bool, trueValue, falseValue T) T {
	if ok {
		return trueValue
	}
	return falseValue
}
