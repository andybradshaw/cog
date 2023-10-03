package loaders

import (
	"fmt"
	"path/filepath"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
)

func kindsysCompopsableLoader(opts Options) ([]*ast.File, error) {
	themaRuntime := thema.NewRuntime(cuecontext.New())

	allSchemas := make([]*ast.File, 0, len(opts.KindsysComposableEntrypoints))
	for _, entrypoint := range opts.KindsysComposableEntrypoints {
		pkg := filepath.Base(entrypoint)

		overlayFS, err := buildKindsysEntrypointFS(opts, entrypoint)
		if err != nil {
			return nil, err
		}

		cueInstance, err := kindsys.BuildInstance(themaRuntime.Context(), ".", "grafanaplugin", overlayFS)
		if err != nil {
			return nil, fmt.Errorf("could not load kindsys instance: %w", err)
		}

		props, err := kindsys.ToKindProps[kindsys.ComposableProperties](cueInstance)
		if err != nil {
			return nil, fmt.Errorf("could not convert cue value to kindsys props: %w", err)
		}

		kindDefinition := kindsys.Def[kindsys.ComposableProperties]{
			V:          cueInstance,
			Properties: props,
		}

		boundKind, err := kindsys.BindComposable(themaRuntime, kindDefinition)
		if err != nil {
			return nil, fmt.Errorf("could not bind kind definition to kind: %w", err)
		}

		schemaAst, err := simplecue.GenerateAST(kindToLatestSchema(boundKind), simplecue.Config{
			Package: pkg, // TODO: extract from input schema/folder?
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}