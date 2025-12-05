package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {

	fmt.Println("FibonacciIterative(10):", FibonacciIterative(10))
	fmt.Println("FibonacciRecursive(10):", FibonacciRecursive(10))

	fmt.Println("IsPrime(2):", IsPrime(2))
	fmt.Println("IsPrime(15):", IsPrime(15))
	fmt.Println("IsPrime(29):", IsPrime(29))

	fmt.Println("IsBinaryPalindrome(7):", IsBinaryPalindrome(7))
	fmt.Println("IsBinaryPalindrome(6):", IsBinaryPalindrome(6))

	fmt.Println(`ValidParentheses("[]{}()"):`, ValidParentheses("[]{}()"))
	fmt.Println(`ValidParentheses("[{]}"):`, ValidParentheses("[{]}"))

	fmt.Println(`Increment("101") ->`, Increment("101"))
	fmt.Println(`Increment("111") ->`, Increment("111"))

}

func FibonacciIterative(n int) int {
	if n < 0 {
		return n
	}
	a, b := 0, 1

	for range n {
		a, b = b, a+b
	}
	return a
}

func FibonacciRecursive(n int) int {
	// switch n {
	// case 0:
	// 	return 0
	// case 1:
	// 	return 1
	// default:
	// 	return fibonacciRecursive(n-1) + fibonacciRecursive(n-2)
	// }
	if n < 2 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

func IsPrime(n int) bool {

	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	limit := int(math.Sqrt(float64(n)))

	for i := 3; i <= limit; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func IsBinaryPalindrome(n int) bool {
	if n < 0 {
		return false
	}

	bin := strconv.FormatInt(int64(n), 2)
	fmt.Println("bin:", bin)

	i, j := 0, len(bin)-1

	for range len(bin) {
		if bin[i] != bin[j] {
			return false
		}
		i, j = i+1, j-1
	}

	return true
}

func ValidParentheses(s string) bool {
	stack := []rune{}
	pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}

	for _, char := range s {
		if pairs[char] == 0 {
			stack = append(stack, char)
		} else if len(stack) == 0 || pairs[char] != stack[len(stack)-1] {
			return false
		} else {
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}
func Increment(num string) int {

	if strings.Trim(num, "01") != "" {
		return 0
	}
	n, err := strconv.ParseInt(num, 2, 64)
	if err != nil {
		return 0
	}
	return int(n + 1)
}
