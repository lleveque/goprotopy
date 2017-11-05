package greet

import "fmt"

//go:generate goprotopy --packagePath=github.com/lleveque/goprotopy/example $GOFILE

// @protopy
func Hello(input []byte) (output []byte, err error) {
    greeting := fmt.Sprintf("Hello %s !", string(input))
    output = []byte(greeting)
    return
}
