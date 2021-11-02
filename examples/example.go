package examples

import (
	"context"
	"github.com/fourhu/go-annotation/pkg/plugin"
	"github.com/mj37yhyy/gowb/pkg/model"
	"github.com/mj37yhyy/gowb/pkg/web"
	"net/http"
)

var _ = plugin.Description{}

//go:generate annotation-gen -i . -v 8

var _ = plugin.Description{}

// Annotation@Description={"body":"a"}
type A struct {
	FieldA string
}

// Annotation@Description={"body":"b"}
type B struct {
	FieldB string
}

// Annotation@Description={"body":"b"}
type C interface {
}

// Annotation@Component
type ComponentA struct {
	B1 *ComponentB `autowired:"true"` // Will populate with new(ComponentB)
	B2 *ComponentB `autowired:"true"` // Will populate with new(ComponentB)
	B3 *ComponentB
}

// Annotation@Component={"type": "Singleton"}
type ComponentB struct {
	C *ComponentC `autowired:"true"` // Will populate with NewComponentC()
}

// Annotation@Component
type ComponentC struct {
	D        *ComponentD `autowired:"true"` // Will populate with NewComponentD()
	IntValue int
}

func NewComponentC() *ComponentC {
	return &ComponentC{IntValue: 1}
}

// Annotation@Component
type ComponentD struct {
	IntValue int
}

func NewComponentD() (*ComponentD, error) {
	return &ComponentD{IntValue: 2}, nil
}

// Annotation@Component
type ApplicationHandlerService struct {
	Order              []string
	ApplicationHandler *ApplicationHandler `autowired:"true"`
}

// Annotation@Service
type ApplicationHandler struct {
}

func ValidateCreateService(ctx context.Context) (model.Response, web.HttpStatus, error) {
	return model.Response{}, http.StatusOK, nil
}

// Annotation@BeforeCallOrder=["Logger"]
// Annotation@AfterCallOrder=["Logger"]
// Annotation@Validate={"enable": true}
func (handler *ApplicationHandler) CreateService(ctx context.Context) (model.Response, web.HttpStatus) {
	return model.Response{}, http.StatusOK
}

func (handler *ApplicationHandler) ModifyService(ctx context.Context) (model.Response, web.HttpStatus) {
	return model.Response{}, http.StatusOK
}

func (handler *ApplicationHandler) DeleteService(ctx context.Context) (model.Response, web.HttpStatus) {
	return model.Response{}, http.StatusOK
}

func (handler *ApplicationHandler) DescribeServices(ctx context.Context) (model.Response, web.HttpStatus) {
	return model.Response{}, http.StatusOK
}

func (handler *ApplicationHandler) CreateServiceVersion(ctx context.Context) (model.Response, web.HttpStatus) {
	return model.Response{}, http.StatusOK
}

func (handler *ApplicationHandler) DescribeServiceVersions(ctx context.Context) (model.Response, web.HttpStatus) {
	return model.Response{}, http.StatusOK
}

func (handler *ApplicationHandler) ModifyServiceVersionReplicas(ctx context.Context) (model.Response, web.HttpStatus) {
	return model.Response{}, http.StatusOK
}

func (handler *ApplicationHandler) DeleteServiceVersion(ctx context.Context) (model.Response, web.HttpStatus) {
	return model.Response{}, http.StatusOK
}
