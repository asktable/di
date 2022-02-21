package di

import (
    "fmt"
    "reflect"
)

var (
    container = map[reflect.Type]map[string]interface{}{}
    listeners = map[reflect.Type]map[string][]reflect.Value{}
)

func Register(ins interface{}) {
    RegisterWithName("", ins)
}

func RegisterWithName(name string, ins interface{}) {
    t := reflect.TypeOf(ins)
    assertPtr(t)
    if _, ok := container[t]; !ok {
        container[t] = map[string]interface{}{}
    }
    container[t][name] = ins
    inject(ins)
    listen(name, ins)
}

func Provide(fn interface{}) {
    if fn == nil {
        panic("fn nil.")
    }
    if reflect.TypeOf(fn).Kind() != reflect.Func {
        panic("fn isn't func.")
    }
    t := reflect.TypeOf(fn)
    args := make([]reflect.Value, 0)
    for i := 0; i < t.NumIn(); i++ {
        in := t.In(i)
        target, ok := find("", in)
        if !ok {
            panic(fmt.Sprintf("not found %v. ", in))
        }
        args = append(args, reflect.ValueOf(target))
    }
    reflect.ValueOf(fn).Call(args)
}

func inject(ins interface{}) {
    t := reflect.TypeOf(ins).Elem()
    v := reflect.ValueOf(ins).Elem()
    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)
        name, ok := f.Tag.Lookup("inject")
        if !ok {
            continue
        }
        target, ok := find(name, f.Type)
        if !ok {
            if listeners[f.Type] == nil {
                listeners[f.Type] = map[string][]reflect.Value{}
            }
            listeners[f.Type][name] = append(listeners[f.Type][name], v.Field(i))
            continue
        }
        set(v.Field(i), target)
    }
}

func listen(name string, ins interface{}) {
    names := []string{""}
    if name != "" {
        names = append(names, name)
    }
    t := reflect.TypeOf(ins)
    target := reflect.ValueOf(ins)
    for _, n := range names {
        for _, v := range listeners[t][n] {
            v.Set(target)
        }
        delete(listeners[t], n)
    }
    if len(listeners[t]) == 0 {
        delete(listeners, t)
    }
}

func find(name string, t reflect.Type) (interface{}, bool) {
    sources, ok := container[t]
    if !ok {
        return nil, false
    }
    if name == "" {
        for _, v := range sources {
            return v, true
        }
    }
    v, ok := sources[name]
    return v, ok
}

func set(obj reflect.Value, target interface{}) {
    obj.Set(reflect.ValueOf(target))
}

func assertPtr(t reflect.Type) {
    if t.Kind() != reflect.Ptr {
        panic(fmt.Sprintf("%v isn't ptr. ", t))
    }
}
