package main

import (
	"fmt"

	"github.com/daviani/go__001/internal/calculator"
)

func main() {
	var name string
	var age int

	name = "roi du go"
	age = 1

	result, isNegative := calculator.Add(5, 3)
	names := "roi du gooo"
	ages := "2"
	fmt.Println("Hello World", name, age, names, ages, result, isNegative)
	fmt.Printf("Type de ages: %T, valeur: %v\n", ages, ages)
	fmt.Printf("Type de age: %T, valeur: %v\n", age, age)

}
