package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/fourhu/go-annotation/pkg/middleware"
	"io"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fourhu/go-annotation/pkg/plugin"
	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/set-gen/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"

	annoArgs "github.com/fourhu/go-annotation/cmd/annotation-gen/args"
	"github.com/fourhu/go-annotation/pkg/lib"
	_ "github.com/fourhu/go-annotation/pkg/plugin"
)

// prefix$Enable=true
// prefix$Type=$Body
type annotation struct {
	rawTypeName string
	body        string
}

// key is annotation rawTypeName
type annotations map[string]*annotation

func isJson(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func annotationEnabled(prefix string, comments []string) bool {
	enableFlag := prefix + "Enable"
	for _, comment := range comments {
		if strings.HasPrefix(comment, enableFlag) {
			return true
		}
	}
	return false
}

func extractAnnotations(prefix string, t *types.Type) annotations {
	ret := map[string]*annotation{}
	comments := append(t.CommentLines, t.SecondClosestCommentLines...)
	for _, comment := range comments {
		var rawTypeName, body string
		if !strings.HasPrefix(comment, prefix) {
			continue
		}
		s := strings.TrimPrefix(comment, prefix)
		sl := strings.Split(s, "=")
		if len(sl) == 1 {
			rawTypeName, body = sl[0], ""
		} else if len(sl) == 2 {
			rawTypeName, body = sl[0], sl[1]
		} else {
			klog.V(4).Infof("annotation format not valid %s\n", comment)
			continue
		}

		if body == "" {
			body = "{}"
		}

		if !isJson(body) {
			klog.V(1).Infoln("annotation format not valid: not valid json", body)
			continue
		}

		ret[rawTypeName] = &annotation{
			rawTypeName: rawTypeName,
			body:        body,
		}
	}
	return ret
}

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPrivateNamer(0, ""),
		"raw":    namer.NewRawNamer("", nil),
	}
}

func DefaultNameSystem() string {
	return "public"
}

// Packages
func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	// LoadGoBoilerplate
	//boilerplate, err := arguments.LoadGoBoilerplate()
	//if err != nil {
	//	klog.Fatalf("Failed loading boilerplate: %v", err)
	//}

	inputs := sets.NewString(context.Inputs...)
	packages := generator.Packages{}
	annotationArgs := arguments.CustomArgs.(*annoArgs.AnnotationArgs)

	// header
	header := append([]byte(fmt.Sprintf("// +build !%s\n\n", arguments.GeneratedBuildTag)) /*boilerplate...*/)

	// arguments handling

	// inputs, get package from context.Universe
	for i := range inputs {
		klog.V(5).Infof("Considering pkg %q", i)
		pkg := context.Universe[i]

		for _, a := range pkg.Imports {
			context.AddDirectory(a.Path)
		}
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}
		//
		if !annotationEnabled(annotationArgs.AnnotationPrefix, pkg.Comments) {
			continue
		}

		klog.V(5).Infof("Generating for pkg %q", i)

		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
				PackagePath: pkg.Path,
				HeaderText:  header,
				// generator 一个 Generator 生成一个文件
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					return []generator.Generator{
						NewGenAnnotation(arguments.OutputFileBaseName, annotationArgs.AnnotationPrefix, pkg),
					}
				},
				// 过滤函数，哪些 type 不关心的，直接过滤，不会被 generator 处理
				// generator 的 过滤器 也可以完成类似的事情，调用时机不同
				FilterFunc: func(c *generator.Context, t *types.Type) bool {
					return t.Name.Package == pkg.Path
				},
			})

	}
	return packages
}

// Order
// 1. Filter()        // Subsequent calls see only types that pass this.
// 2. Namers()        // Subsequent calls see the namers provided by this.
// 3. PackageVars()
// 4. PackageConsts()
// 5. Init()
// 6. GenerateType()  // Called N times, once per type in the context's Order.
// 7. Imports()
type genAnnotation struct {
	generator.DefaultGen
	targetPackage    string
	annotationPrefix string
	pkg              *types.Package
	imports          namer.ImportTracker
	importsCache     []string
}

func NewGenAnnotation(sanitizedName, annotationPrefix string, pkg *types.Package) generator.Generator {
	return &genAnnotation{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		pkg:              pkg,
		targetPackage:    pkg.Path,
		annotationPrefix: annotationPrefix,
		imports:          generator.NewImportTracker(),
		importsCache: []string{
			"context",
			"github.com/fourhu/go-annotation/pkg/lib",
			"k8s.io/klog",
			"github.com/mj37yhyy/gowb/pkg/web",
			"github.com/mj37yhyy/gowb/pkg/model",
			"github.com/fourhu/go-annotation/pkg/middleware",
		},
	}
}

// Namer for template
func (g *genAnnotation) Namers(c *generator.Context) namer.NameSystems {
	// Have the raw namer for this file track what it imports.
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.targetPackage, g.imports),
	}
}

//
func (g *genAnnotation) Filter(c *generator.Context, t *types.Type) bool {
	return true
}

//
func (g *genAnnotation) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, g.importsCache...)
	return
}

// Init method for generated code
func (g *genAnnotation) Init(c *generator.Context, w io.Writer) error {
	return nil
}

func findAnnotationType(c *generator.Context, name string) *types.Type {
	for _, p := range c.Universe {
		for _, t := range p.Types {
			klog.V(8).Infoln("finding name", t.Name.Name)
			if t.Name.Name == name {
				return t
			}
		}
	}
	return nil
}

func (g *genAnnotation) getOriginServiceName(t *types.Type) string {
	return t.Name.Name
}

func (g *genAnnotation) getNewFunction(c *generator.Context, t *types.Type) string {
	for name, method := range g.pkg.Functions {
		if name == "New"+g.Namers(c)["raw"].Name(t) {
			// do not support parameters now
			if method.Underlying == nil || method.Underlying.Signature == nil {
				continue
			}
			signature := method.Underlying.Signature

			signatureMatch := func(p *types.Type, s *types.Type) bool {
				return g.Namers(c)["raw"].Name(p) == "*"+s.Name.Name
			}

			if len(signature.Parameters) == 0 {
				if len(signature.Results) > 0 {
					fmt.Println(g.Namers(c)["raw"].Name(signature.Results[0]), "*", t.Name.Name)
				}

				if len(signature.Results) > 1 {
					fmt.Println(g.Namers(c)["raw"].Name(signature.Results[1]), "*", t.Name.Name)
				}

				if len(signature.Results) == 1 &&
					signatureMatch(signature.Results[0], t) {
					return fmt.Sprintf(`func () (interface{}, error) { return %s(), nil }`, name)
				}
				if len(signature.Results) == 2 &&
					signatureMatch(signature.Results[0], t) &&
					signature.Results[1].Name.Name == "error" {
					return fmt.Sprintf(`func () (interface{}, error) { return %s() }`, name)
				}
			}
		}
	}
	return fmt.Sprintf(`func () (interface{}, error) { return new(%s), nil}`, g.Namers(c)["raw"].Name(t))
}

func (g *genAnnotation) trimImport(name string) string {
	for _, imt := range g.importsCache {
		if strings.HasPrefix(name, imt) {
			i := path.Base(strings.ReplaceAll(imt, fmt.Sprintf(".%s", name), ""))
			if i != "" {
				return fmt.Sprintf("%s.%s", i, strings.ReplaceAll(name, fmt.Sprintf("%s.", imt), ""))
			}
		}
	}
	return name
}

// core
func (g *genAnnotation) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	// writer
	//sw := generator.NewSnippetWriter(w, c, "$", "$")
	sw := generator.NewSnippetWriter(w, c, "%%", "%%")

	klog.V(5).Infof("processing type %v", t)
	// params
	annotations := extractAnnotations(g.annotationPrefix, t)

	for _, anno := range annotations {
		annotationType := findAnnotationType(c, anno.rawTypeName)
		if annotationType == nil {
			klog.V(1).Infoln("annotation type not found", anno.rawTypeName)
			continue
		}
		annotationPlugin := lib.GetPluginByName(anno.rawTypeName)
		var annotationTarget = lib.TargetDefault
		if targeted, ok := annotationPlugin.(lib.TargetedAnnotation); ok {
			annotationTarget = targeted.Target()
		}

		if annotationTarget != t.Kind {
			continue
		}

		pluginPopulated := reflect.New(reflect.TypeOf(annotationPlugin).Elem()).Interface()
		err := json.Unmarshal([]byte(anno.body), pluginPopulated)
		if anno.body != "" && err != nil {
			klog.V(1).Infoln("annotation unmarshal error", err)
			continue
		}
		newFunctionIsSingleton := false

		klog.V(8).Infoln("pluginPopulated ", reflect.TypeOf(annotationPlugin), reflect.TypeOf(pluginPopulated))
		if p, ok := pluginPopulated.(*plugin.Component); ok {

			newFunctionIsSingleton = p.Type == plugin.Singleton
		}

		klog.V(1).Infoln("getNewFunction", g.getNewFunction(c, t))

		serviceMethodMap := g.getServiceMethodMap(t)
		m := map[string]interface{}{
			"Resource":             c.Universe.Function(types.Name{Package: t.Name.Package, Name: "Resource"}),
			"type":                 t,
			"annotationType":       annotationType,
			"annotationBody":       anno.body,
			"newFunction":          g.getNewFunction(c, t),
			"getOriginServiceName": g.getOriginServiceName(t),
			"newFunctionSingleton": newFunctionIsSingleton,
			"ServiceMethodMap":     serviceMethodMap,
		}
		klog.V(3).Infoln("annotation m", m)

		if anno.rawTypeName != "Service" {
			// render registerTemplate
			sw.Do(registerTemplate, m)
		}

		if compile, ok := annotationPlugin.(lib.CompileAnnotation); ok {
			// render plugin template
			sw.Do(compile.Template(), m)
		} else {
			klog.V(4).Infoln("get compile", annotationPlugin)
		}
	}

	return sw.Error()
}

type methodInfo struct {
	Name           string `json:"Name"`
	EnableValidate bool   `json:"EnableValidate"`

	Parameters     []string `json:"Parameters"`
	ParameterTypes []string `json:"ParameterTypes"`
	ParameterExpr  string   `json:"ParameterExpress"`
	ResultTypes    []string `json:"ResultTypes"`
	Results        string   `json:"Results"`
	ResultVarExpr  string   `json:"ResultVarExpr"`
	ResultExpr     string   `json:"ResultExpr"`
}

func (g *genAnnotation) getServiceMethodMap(t *types.Type) map[string]*methodInfo {
	serviceMethodMap := make(map[string]*methodInfo, 5)
	for name, tm := range t.Methods {
		m := &methodInfo{Name: name}
		var parameters []string
		var parameterTypes []string
		var parameterExpress []string

		// params
		annotations := extractAnnotations(g.annotationPrefix, tm)
		if len(annotations) == 0 {
			continue
		}
		for _, anno := range annotations {
			if anno.rawTypeName != "Expose" {
				continue
			}

			validateEnable := new(ValidateEnable)
			err := json.Unmarshal([]byte(anno.body), validateEnable)
			if err != nil {
				klog.V(4).Infof("json unmarshal validate err: %v", err)
				continue
			}

			if validateEnable.Enable {
				m.EnableValidate = true
			}
		}

		parameters = append(parameters, "ctx")
		parameterTypes = append(parameters, "context.Context")
		parameterExpress = append(parameterExpress, "ctx context.Context")
		for idx, p := range tm.Signature.Parameters {
			t := fmt.Sprintf("%s", g.trimImport(p.Name.String()))
			if strings.HasPrefix(t, "context.Context") {
				continue
			}
			v := fmt.Sprintf("arg%d", idx)
			parameterExpress = append(parameters, fmt.Sprintf("%s %s", v, t))
			parameters = append(parameters, v)
			parameterTypes = append(parameterTypes, t)
		}

		m.ParameterTypes = parameterTypes
		m.Parameters = parameters

		if len(parameters) > 0 {
			m.ParameterExpr = fmt.Sprintf("(%s)", strings.Join(parameterExpress, ","))
		} else {
			m.ParameterExpr = "()"
		}

		var resultTypes []string
		var results []string
		var resultVarExpress []string
		var resultExpress []string
		for idx, r := range tm.Signature.Results {
			t := g.trimImport(r.Name.String())
			v := fmt.Sprintf("ret%d", idx)
			resultTypes = append(resultTypes, t)
			results = append(results, v)
			resultExpress = append(resultExpress, fmt.Sprintf("%s", v))
			resultVarExpress = append(resultVarExpress, fmt.Sprintf("%s %s", v, t))
		}

		m.ResultTypes = resultTypes

		if len(resultTypes) == 1 {
			m.Results = fmt.Sprintf("%s", resultTypes[0])
		} else if len(resultTypes) > 1 {
			m.Results = fmt.Sprintf("(%s)", strings.Join(resultTypes, ","))
		}

		if m.EnableValidate {
			if !middleware.StrInSlice("err error", resultVarExpress) {
				resultVarExpress = append(resultVarExpress, "err error")
			}
		}

		if len(resultVarExpress) == 1 {
			m.ResultVarExpr = fmt.Sprintf("var %s", resultVarExpress[0])
		} else if len(resultVarExpress) > 1 {
			m.ResultVarExpr = fmt.Sprintf("var (\n%s\n)", strings.Join(resultVarExpress, "\n"))
		} else {
			m.ResultVarExpr = ""
		}

		m.ResultExpr = fmt.Sprintf("%s", strings.Join(resultExpress, ","))
		serviceMethodMap[name] = m
	}

	klog.V(3).Infof("ServiceMethodMap: %+v", serviceMethodMap)
	return serviceMethodMap
}

type ValidateEnable struct {
	Enable bool `json:"enableValidate"`
}

// register template
var registerTemplate = `
func init() {
	b := new(%%.annotationType|raw%%)
	err := json.Unmarshal([]byte("%%.annotationBody|js%%"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new(%%.type|raw%%), b)
}

`
