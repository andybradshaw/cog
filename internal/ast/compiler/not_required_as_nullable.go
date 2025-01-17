package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*NotRequiredFieldAsNullableType)(nil)

// NotRequiredFieldAsNullableType identifies all the struct fields marked as not `Required`
// and rewrites their `Type` to be `Nullable`.
type NotRequiredFieldAsNullableType struct {
}

func (pass *NotRequiredFieldAsNullableType) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *NotRequiredFieldAsNullableType) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		return pass.processObject(object)
	})

	return schema
}

func (pass *NotRequiredFieldAsNullableType) processObject(object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	object.Type = pass.processType(object.Type)

	return object
}

func (pass *NotRequiredFieldAsNullableType) processType(def ast.Type) ast.Type {
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

func (pass *NotRequiredFieldAsNullableType) processArray(def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(def.Array.ValueType)

	return def
}

func (pass *NotRequiredFieldAsNullableType) processMap(def ast.Type) ast.Type {
	def.Map.IndexType = pass.processType(def.Map.IndexType)
	def.Map.ValueType = pass.processType(def.Map.ValueType)

	return def
}

func (pass *NotRequiredFieldAsNullableType) processDisjunction(def ast.Type) ast.Type {
	for i, branch := range def.Disjunction.Branches {
		def.Disjunction.Branches[i] = pass.processType(branch)
	}

	return def
}

func (pass *NotRequiredFieldAsNullableType) processStruct(def ast.Type) ast.Type {
	for i, field := range def.Struct.Fields {
		def.Struct.Fields[i].Type = pass.processType(field.Type)
		if !field.Required && !def.Struct.Fields[i].Type.Nullable {
			def.Struct.Fields[i].Type.Nullable = true
			def.Struct.Fields[i].AddToPassesTrail("NotRequiredFieldAsNullableType[nullable=true]")
		}
	}

	return def
}
