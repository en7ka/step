package main

import "fmt"

func main() {
	list := []int{1, 2, 3, 4, 5, 6}

	fmt.Println(double(list))
}

func double(nums []int) []int {
	var res []int

	for _, num := range nums {
		res = append(res, num*2)
	}
	return res
}
