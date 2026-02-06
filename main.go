package main

import (
	"errors"
	"fmt"
)

func main() {
	var name string
	var age int

	name = "roi du go"
	age = 1

	result, isNegative := add(5, 3)
	names := "roi du gooo"
	ages := "2"
	fmt.Println("Hello World", name, age, names, ages, result, isNegative)
	fmt.Printf("Type de ages: %T, valeur: %v\n", ages, ages)
	fmt.Printf("Type de age: %T, valeur: %v\n", age, age)

}

func add(a, b int) (int, error) {
	result := a + b
	if result < 0 {
		return 0, errors.New("retour négatif")
	}
	return result, nil
}
