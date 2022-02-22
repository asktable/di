# Introduce
Dependency injection for Go.
# Getting Started
Install
```shell
go get -u github.com/asktable/di
```
Usage
```golang
package main

import (
    "fmt"
    "github.com/asktable/di"
)

type Car struct {
    Engine *Engine `inject:""`
}

type Engine struct {
    Name string
}

func (e *Engine) Run() {
    fmt.Println(e.Name, "is running")
}

func main() {
    di.Register(&Car{})
    di.Register(&Engine{Name: "DI"})
    di.Provide(func(car *Car) {
        car.Engine.Run()    // print: DI is running 
    })
}
```