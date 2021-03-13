package day2

import (
    "fmt"

    "example.com/helloworld"
)

func main() {
    // Get a greeting message and print it.
    const name string = "Gladys"
    message := helloworld.Hello(name)
    fmt.Println(message)
    helloworld.Bye(name)
}
