package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DisjunctionWithNullToOptional)(nil)

// DisjunctionWithNullToOptional simplifies disjunctions with two branches, where one is `null`. For those,
// it transforms `type | null` into `*type` (optional, nullable reference to `type`).
//
// Example:
//
//	```
//	MaybeString: string | null
//	```
//
// Will become:
//
//	```
//	MaybeString?: string
//	```
type DisjunctionWithNullToOptional struct {
}

func (pass *DisjunctionWithNullToOptional) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		newSchema, err := pass.processSchema(schema)
		if err != nil {
			return nil, fmt.Errorf("[%s] %w", schema.Package, err)
		}

		newSchemas = append(newSchemas, newSchema)
	}

	return newSchemas, nil
}

func (pass *DisjunctionWithNullToOptional) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		return pass.processObject(object)
	})

	return schema, nil
}

func (pass *DisjunctionWithNullToOptional) processObject(object ast.Object) ast.Object {
	object.Type = pass.processType(object.Type)

	return object
}

func (pass *DisjunctionWithNullToOptional) processType(def ast.Type) ast.Type {
	if def.IsArray() {
		return pass.processArray(def)
	}

	if def.IsMap() {
		return pass.processMap(def)
	}

	if def.IsStruct() {
		return pass.processStruct(def)
	}

	if def.IsDisjunction() {
		return pass.processDisjunction(def)
	}

	return def
}

func (pass *DisjunctionWithNullToOptional) processArray(def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(def.AsArray().ValueType)

	return def
}

func (pass *DisjunctionWithNullToOptional) processMap(def ast.Type) ast.Type {
	def.Map.ValueType = pass.processType(def.AsMap().ValueType)

	return def
}

func (pass *DisjunctionWithNullToOptional) processStruct(def ast.Type) ast.Type {
	for i, field := range def.AsStruct().Fields {
		def.Struct.Fields[i].Type = pass.processType(field.Type)
	}

	return def
}

func (pass *DisjunctionWithNullToOptional) processDisjunction(def ast.Type) ast.Type {
	disjunction := def.AsDisjunction()

	// type | null
	if len(disjunction.Branches) == 2 && disjunction.Branches.HasNullType() {
		finalType := disjunction.Branches.NonNullTypes()[0]
		finalType.Nullable = true
		finalType.AddToPassesTrail(fmt.Sprintf("DisjunctionWithNullToOptional[%[1]s|null → %[1]s?]", ast.TypeName(finalType)))

		return finalType
	}

	return def
}
