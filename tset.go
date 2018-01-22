package main

import (
    "fmt"
)

type Test struct {
    num     int
    msg     *string
}

func NewTestDef() (*Test, error) {
    msg := "It's a default"
    t, e := NewTest(&msg)

    return t, e
}

func NewTest(msg *string) (*Test, error) {
    test := &Test {
        num: 0,
        msg: msg,
    }

    return test, nil
}

func (t *Test) exec(endline *string) error {
    fmt.Println(*t.msg + " " + *endline)
    return nil
}

func main() {
    test, _ := NewTestDef()
    name := "Taro"
    test.exec(&name)
}

