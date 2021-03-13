package main

import (
    "fmt"
    "strings"

    "example.com/helloworld"
)

func main() {
    // Get a greeting message and print it.
    const name string = "Lunit"
    var message = helloworld.Hello(name)
    fmt.Println(message)
    nameUpper := strings.ToUpper(name)
    fmt.Println(nameUpper)
    messageUpper := strings.ToUpper(message)
    fmt.Println(messageUpper)
    helloworld.Bye(name)
}
