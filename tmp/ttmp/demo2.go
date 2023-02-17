package main

import "fmt"

func main() {
	// 在控制台上输出带颜色的文本
	fmt.Println("\033[31mThis text is red.\033[0m")
	fmt.Println("\033[32mThis text is green.\033[0m")
	fmt.Println("\033[33mThis text is yellow.\033[0m")
	fmt.Println("\033[34mThis text is blue.\033[0m")
	fmt.Println("\033[35mThis text is magenta.\033[0m")
	fmt.Println("\033[36mThis text is cyan.\033[0m")
}
