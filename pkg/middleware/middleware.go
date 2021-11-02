package middleware

import (
	"context"
	"sync"
)

var (
	CallListMap *sync.Map
)

func init() {
	CallListMap = new(sync.Map)
	RegisterCall("Logger", &Logger{})
}

func RegisterCall(name string, call CallList) {
	CallListMap.Store(name, call)
}

type CallList interface {
	Handler(c context.Context, args ...interface{})
}

func Before(c context.Context, order []string, args ...interface{}) {
	callMap := make(map[string]CallList, 10)
	CallListMap.Range(func(key, value interface{}) bool {
		name := key.(string)
		if StrInSlice(name, order) {
			f := value.(CallList)
			callMap[name] = f
		}

		return false
	})

	for _, o := range order {
		f, ok := callMap[o]
		if ok {
			f.Handler(c, args)
		}
	}
}

func StrInSlice(str string, args []string) bool {
	if len(args) == 0 {
		return false
	}

	for _, s := range args {
		if s == str {
			return true
		}
	}

	return false
}

func After(c context.Context, order []string, args ...interface{}){
	callMap := make(map[string]CallList, 10)
	CallListMap.Range(func(key, value interface{}) bool {
		name := key.(string)
		if StrInSlice(name, order) {
			f := value.(CallList)
			callMap[name] = f
		}

		return false
	})

	for _, o := range order {
		f, ok := callMap[o]
		if ok {
			f.Handler(c, args)
		}
	}
}


