package option

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type Selector func(builder ast.Builder, option ast.Option) bool

func ByName(pkg string, objectName string, optionNames ...string) Selector {
	return func(builder ast.Builder, option ast.Option) bool {
		return builder.For.SelfRef.ReferredPkg == pkg &&
			builder.For.Name == objectName &&
			tools.ItemInList(option.Name, optionNames)
	}
}

func ByNameCaseInsensitive(pkg string, objectName string, optionNames ...string) Selector {
	return func(builder ast.Builder, option ast.Option) bool {
		return builder.For.SelfRef.ReferredPkg == pkg &&
			builder.For.Name == objectName &&
			tools.StringInListEqualFold(option.Name, optionNames)
	}
}
