package di

import (
    "fmt"
    "testing"
)

type Api struct {
    S1 *Service `inject:"s1"`
    S2 *Service `inject:"s2"`
    DB *DB      `inject:""`
}

type Service struct {
    Name string
}

type DB struct {
    Name string
}

func TestRegister(t *testing.T) {
    Register(&Api{})
    Register(&DB{"db"})
    RegisterWithName("s1", &Service{"s1"})
    RegisterWithName("s2", &Service{"s2"})

    Provide(func(a *Api, s *Service) {
        fmt.Println(a.S1.Name)
        fmt.Println(a.S2.Name)
        fmt.Println(a.DB.Name)
    })
}
