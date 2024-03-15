package main

import (
	"flag"
	"fmt"
)

func main() {
	// Определение флагов командной строки
	var name string
	flag.StringVar(&name, "name", "World", "a name to say hello to")

	// Парсинг флагов
	flag.Parse()

	// Вывод аргументов командной строки
	fmt.Printf("Hello, %s!\n", name)

	// Вывод оставшихся аргументов после флагов
	fmt.Println("Other args:", flag.Args())
}
