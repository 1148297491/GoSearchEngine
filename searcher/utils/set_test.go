package utils

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	var arr = []int{1, 2, 3}
	fuc(arr)
	fmt.Printf("%+v", arr)
}

func fuc(arr []int) {
	arr[0] = 7
	arr = append(arr, 5)
}
