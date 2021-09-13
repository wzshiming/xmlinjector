package xmlinjector

import (
	"os"
	"reflect"
	"strings"
)

func NewArgs(tag string, env bool) *Args {
	tag = strings.ReplaceAll(tag, "\n", " ")
	return &Args{
		tag: reflect.StructTag(tag),
		env: env,
	}
}

type Args struct {
	tag reflect.StructTag
	env bool
}

func (a Args) Get(name string) (string, bool) {
	val, ok := a.tag.Lookup(name)
	if !ok {
		return "", false
	}
	if a.env {
		val = os.Expand(val, os.Getenv)
	}
	return val, true
}
