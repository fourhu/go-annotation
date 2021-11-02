package plugin

import (
	"github.com/fourhu/go-annotation/pkg/lib"
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
