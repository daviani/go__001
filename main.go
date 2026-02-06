package main

import (
	"fmt"

	"github.com/daviani/go__001/internal/calculator"
	"github.com/daviani/go__001/internal/user"
)

func main() {
	var name string
	var age int

	name = "roi du go"
	age = 1

	myUser := user.User{Name: "me", Age: 12}
	result, isNegative := calculator.Add(5, 3)

	names := "roi du gooo"
	ages := "2"
	fmt.Println("Hello World", name, age, names, ages, result, isNegative, myUser.Greet())
	fmt.Printf("Type de ages: %T, valeur: %v\n", ages, ages)
	fmt.Printf("Type de age: %T, valeur: %v\n", age, age)

}
