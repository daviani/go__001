package main

import "fmt"

func main() {
	var name string
	var age int

	name = "roi du go"
	age = 1

	names := "roi du gooo"
	ages := "2"
	fmt.Println("Hello World", name, age, names, ages)
	fmt.Printf("Type de ages: %T, valeur: %v\n", ages, ages)
	fmt.Printf("Type de age: %T, valeur: %v\n", age, age)

}
