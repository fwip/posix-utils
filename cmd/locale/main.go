package main

import "fmt"
import "github.com/fwip/posix-utils/pkg/locale"

func main() {
	fmt.Printf(locale.FromEnv().String())
}
