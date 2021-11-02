package plugin

import (
	"github.com/u2takey/go-annotation/pkg/lib"
)

func init() {
	lib.RegisterPlugin(new(Service))
}

//type ServiceType string

//const (
//Default   ServiceType = ""
//Singleton ServiceType = "Singleton"
//)

type Service struct {
}

func (p *Service) Template() string {
	// register a New Method
	return `
type %%.type|raw%%Service struct {
	Order []string
	%%.type|raw%% *%%.type|raw%% ` + "`autowired:\"true\"`" +
		`}

var New%%.type|raw%%ServiceFunction = &lib.NewFunction{
	F: %%.newFunction%%,
	Singleton: %%.newFunctionSingleton|print%%,
}

func init() {
	lib.RegisterType(new(%%.type|raw%%Service), New%%.type|raw%%ServiceFunction) 
}

func Provide%%.type|raw%%Service () (*%%.type|raw%%Service, error) {
	r, err := lib.Provide(new(%%.type|raw%%Service))
	if err != nil{
		return nil, err
	}
	return r.(*%%.type|raw%%Service), nil
}

%%$handlerType := .type|raw%%
%%range $name, $info := .ServiceMethodMap%%
func (handler *%%$handlerType%%Service) %%$name%% %%$info.ParameterExpr%% %%$info.Results%% {
	%%$info.ResultVarExpr%%
    middleware.Before(ctx, handler.Order)
	%% if $info.EnableValidate%%
	%%$info.ResultExpr%%, err = Validate%%$name%%(ctx)
	if err != nil {
		return %%$info.ResultExpr%%
	}
	%%else%%
	%%end%%

	%%$info.ResultExpr%% = handler.ApplicationHandler.CreateService(ctx)
	middleware.After(ctx, handler.Order, %%$info.ResultExpr%%)
	return %%$info.ResultExpr%%
}
%%end%%
`
}
