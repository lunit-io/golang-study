package main

import (
    "fmt"

    "example.com/helloworld"
)

func main() {
    // Get a greeting message and print it.
    const name string = "Lunit"
    message := helloworld.Hello(name)
    fmt.Println(message)
    helloworld.Bye(name)
}
