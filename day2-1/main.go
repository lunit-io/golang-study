package main

import (
    "byebye"
    "fmt"
    "strings"

    "example.com/helloworld"
    "example.com/uuid"
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
    getUuid()
    helloworld.Bye(name)
    fmt.Println(bye.Bye("scbang"))
}


func getUuid() {
    uuidString := uuid.GenerateUUID()
    fmt.Println(uuidString)
}
