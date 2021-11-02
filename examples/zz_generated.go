//go:build !ignore_autogenerated
// +build !ignore_autogenerated

package examples

import (
	"context"
	"encoding/json"

	"github.com/mj37yhyy/gowb/pkg/model"
	"github.com/mj37yhyy/gowb/pkg/web"
	"github.com/u2takey/go-annotation/pkg/lib"
	"github.com/u2takey/go-annotation/pkg/middleware"
	plugin "github.com/u2takey/go-annotation/pkg/plugin"
	"k8s.io/klog"
)

func init() {
	b := new(plugin.Description)
	err := json.Unmarshal([]byte("{\"body\":\"a\"}"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new(A), b)
}

func (s *A) GetDescription() string {
	return "A"
}

func init() {
	b := new(plugin.Component)
	err := json.Unmarshal([]byte("{}"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new(ApplicationHandler), b)
}

var NewApplicationHandlerFunction = &lib.NewFunction{
	F:         func() (interface{}, error) { return new(ApplicationHandler), nil },
	Singleton: false,
}

func init() {
	lib.RegisterType(new(ApplicationHandler), NewApplicationHandlerFunction)
}

func ProvideApplicationHandler() (*ApplicationHandler, error) {
	r, err := lib.Provide(new(ApplicationHandler))
	if err != nil {
		return nil, err
	}
	return r.(*ApplicationHandler), nil
}

func init() {
	b := new(plugin.Service)
	err := json.Unmarshal([]byte("{}"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new(ApplicationHandler), b)
}

type ApplicationHandlerService struct {
	Order              []string
	ApplicationHandler *ApplicationHandler `autowired:"true"`
}

var NewApplicationHandlerServiceFunction = &lib.NewFunction{
	F:         func() (interface{}, error) { return new(ApplicationHandler), nil },
	Singleton: false,
}

func init() {
	lib.RegisterType(new(ApplicationHandlerService), NewApplicationHandlerServiceFunction)
}

func ProvideApplicationHandlerService() (*ApplicationHandlerService, error) {
	r, err := lib.Provide(new(ApplicationHandlerService))
	if err != nil {
		return nil, err
	}
	return r.(*ApplicationHandlerService), nil
}

func (handler *ApplicationHandlerService) CreateService(ctx context.Context) (model.Response, web.HttpStatus) {
	var (
		ret0 model.Response
		ret1 web.HttpStatus
		err  error
	)
	middleware.Before(ctx, handler.Order)

	ret0, ret1, err = ValidateCreateService(ctx)
	if err != nil {
		return ret0, ret1
	}

	ret0, ret1 = handler.ApplicationHandler.CreateService(ctx)
	middleware.After(ctx, handler.Order, ret0, ret1)
	return ret0, ret1
}

func (handler *ApplicationHandlerService) CreateServiceVersion(ctx context.Context) (model.Response, web.HttpStatus) {
	var (
		ret0 model.Response
		ret1 web.HttpStatus
	)
	middleware.Before(ctx, handler.Order)

	ret0, ret1 = handler.ApplicationHandler.CreateService(ctx)
	middleware.After(ctx, handler.Order, ret0, ret1)
	return ret0, ret1
}

func (handler *ApplicationHandlerService) DeleteService(ctx context.Context) (model.Response, web.HttpStatus) {
	var (
		ret0 model.Response
		ret1 web.HttpStatus
	)
	middleware.Before(ctx, handler.Order)

	ret0, ret1 = handler.ApplicationHandler.CreateService(ctx)
	middleware.After(ctx, handler.Order, ret0, ret1)
	return ret0, ret1
}

func (handler *ApplicationHandlerService) DeleteServiceVersion(ctx context.Context) (model.Response, web.HttpStatus) {
	var (
		ret0 model.Response
		ret1 web.HttpStatus
	)
	middleware.Before(ctx, handler.Order)

	ret0, ret1 = handler.ApplicationHandler.CreateService(ctx)
	middleware.After(ctx, handler.Order, ret0, ret1)
	return ret0, ret1
}

func (handler *ApplicationHandlerService) DescribeServiceVersions(ctx context.Context) (model.Response, web.HttpStatus) {
	var (
		ret0 model.Response
		ret1 web.HttpStatus
	)
	middleware.Before(ctx, handler.Order)

	ret0, ret1 = handler.ApplicationHandler.CreateService(ctx)
	middleware.After(ctx, handler.Order, ret0, ret1)
	return ret0, ret1
}

func (handler *ApplicationHandlerService) DescribeServices(ctx context.Context) (model.Response, web.HttpStatus) {
	var (
		ret0 model.Response
		ret1 web.HttpStatus
	)
	middleware.Before(ctx, handler.Order)

	ret0, ret1 = handler.ApplicationHandler.CreateService(ctx)
	middleware.After(ctx, handler.Order, ret0, ret1)
	return ret0, ret1
}

func (handler *ApplicationHandlerService) ModifyService(ctx context.Context) (model.Response, web.HttpStatus) {
	var (
		ret0 model.Response
		ret1 web.HttpStatus
	)
	middleware.Before(ctx, handler.Order)

	ret0, ret1 = handler.ApplicationHandler.CreateService(ctx)
	middleware.After(ctx, handler.Order, ret0, ret1)
	return ret0, ret1
}

func (handler *ApplicationHandlerService) ModifyServiceVersionReplicas(ctx context.Context) (model.Response, web.HttpStatus) {
	var (
		ret0 model.Response
		ret1 web.HttpStatus
	)
	middleware.Before(ctx, handler.Order)

	ret0, ret1 = handler.ApplicationHandler.CreateService(ctx)
	middleware.After(ctx, handler.Order, ret0, ret1)
	return ret0, ret1
}

func init() {
	b := new(plugin.Description)
	err := json.Unmarshal([]byte("{\"body\":\"b\"}"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new(B), b)
}

func (s *B) GetDescription() string {
	return "B"
}

func init() {
	b := new(plugin.Component)
	err := json.Unmarshal([]byte("{}"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new(ComponentA), b)
}

var NewComponentAFunction = &lib.NewFunction{
	F:         func() (interface{}, error) { return new(ComponentA), nil },
	Singleton: false,
}

func init() {
	lib.RegisterType(new(ComponentA), NewComponentAFunction)
}

func ProvideComponentA() (*ComponentA, error) {
	r, err := lib.Provide(new(ComponentA))
	if err != nil {
		return nil, err
	}
	return r.(*ComponentA), nil
}

func init() {
	b := new(plugin.Component)
	err := json.Unmarshal([]byte("{\"type\": \"Singleton\"}"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new(ComponentB), b)
}

var NewComponentBFunction = &lib.NewFunction{
	F:         func() (interface{}, error) { return new(ComponentB), nil },
	Singleton: true,
}

func init() {
	lib.RegisterType(new(ComponentB), NewComponentBFunction)
}

func ProvideComponentB() (*ComponentB, error) {
	r, err := lib.Provide(new(ComponentB))
	if err != nil {
		return nil, err
	}
	return r.(*ComponentB), nil
}

func init() {
	b := new(plugin.Component)
	err := json.Unmarshal([]byte("{}"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new(ComponentC), b)
}

var NewComponentCFunction = &lib.NewFunction{
	F:         func() (interface{}, error) { return NewComponentC(), nil },
	Singleton: false,
}

func init() {
	lib.RegisterType(new(ComponentC), NewComponentCFunction)
}

func ProvideComponentC() (*ComponentC, error) {
	r, err := lib.Provide(new(ComponentC))
	if err != nil {
		return nil, err
	}
	return r.(*ComponentC), nil
}

func init() {
	b := new(plugin.Component)
	err := json.Unmarshal([]byte("{}"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new(ComponentD), b)
}

var NewComponentDFunction = &lib.NewFunction{
	F:         func() (interface{}, error) { return NewComponentD() },
	Singleton: false,
}

func init() {
	lib.RegisterType(new(ComponentD), NewComponentDFunction)
}

func ProvideComponentD() (*ComponentD, error) {
	r, err := lib.Provide(new(ComponentD))
	if err != nil {
		return nil, err
	}
	return r.(*ComponentD), nil
}
