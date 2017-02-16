# failer

Extended format of errors for go-lang.

Features:
 * code number;
 * error tagging;
 * decoding/encoding to bytes;

TODO:
 * support stack catch;


Exmaple usage:


```go
package main

import (
    "github.com/Kerah/failer"
    "fmt"
    "log"
)


func Func1() error {
    errWithCode := failer.New("error message", 42)
    if errWithCode.Code() != 42 {
        panic("unexpected code number")
    }
    return errWithCode
}

func Func2() ([]byte) {
    res := Func1()
    if res != nil {
        return res.(failer.Encoder).Encode()
    }
    return nil
}

func main(){
    data := Func2()
    if data != nil {
        fail, err := failer.Decode(data)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(fail.Error())
    }
}

```