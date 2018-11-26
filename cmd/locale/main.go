package main

import "fmt"
import "github.com/fwip/posix-utils/locale"

func main() {
	fmt.Printf(locale.FromEnv().String())
}
