// +build !ignore_autogenerated

package examples

import (
	"encoding/json"

	"github.com/u2takey/go-annotation/pkg/lib"
	plugins "github.com/u2takey/go-annotation/pkg/plugins"
	"k8s.io/klog"
)

func init() {
	b := new(plugins.Description)
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
	b := new(plugins.Description)
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
