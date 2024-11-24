/*
@author: sk
@date: 2024/11/16
*/
package main

import "fmt"

// https://readonly.link/books/https://raw.githubusercontent.com/GuoYaxiang/craftinginterpreters_zh/main/book.json/-/14.%E5%AD%97%E8%8A%82%E7%A0%81%E5%9D%97.md

func main() {
	source := ReadFile("res/test.txt")
	chunk := Compile(source)
	fmt.Println("=====================Disassemble=====================")
	Disassemble(chunk)
	fmt.Println("=====================Interpret=====================")
	Interpret(chunk)
}
