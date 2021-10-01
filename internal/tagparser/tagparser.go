package tagparser

import (
	"context"
	"errors"
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"edholm.dev/go-logging"
)

type AvailableConfig struct {
	Name    string
	Package string
	Options []ConfigOption
}

func Analyze(ctx context.Context, files []*ast.File) []AvailableConfig {
	availConfigs := make([]AvailableConfig, 0, len(files))

	for _, file := range files {
		ast.Inspect(file, func(node ast.Node) bool {
			if gendecl, ok := node.(*ast.GenDecl); ok {
				if gendecl.Tok != token.TYPE {
					return true
				}

				// TODO: check for multiple
				if len(gendecl.Specs) != 1 {
					return true
				}

				typespec, ok := gendecl.Specs[0].(*ast.TypeSpec)
				if !ok {
					return true
				}

				strukt, ok := typespec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				if !hasTags(strukt) {
					return true
				}

				options := extractConfigOptions(ctx, strukt)
				if len(options) > 0 {
					availConfigs = append(availConfigs, AvailableConfig{
						Name:    typespec.Name.Name,
						Package: file.Name.Name,
						Options: options,
					})
				}
			}

			return true
		})
	}
	return availConfigs
}

func hasTags(strukt *ast.StructType) bool {
	for _, field := range strukt.Fields.List {
		if field.Tag != nil {
			return true
		}
	}
	return false
}

func extractConfigOptions(ctx context.Context, strukt *ast.StructType) []ConfigOption {
	logger := logging.FromContext(ctx)
	options := make([]ConfigOption, 0, len(strukt.Fields.List))

	for _, field := range strukt.Fields.List {
		if field.Tag == nil {
			continue
		}

		rawTag := strings.ReplaceAll(field.Tag.Value, "`", "")
		tag := reflect.StructTag(rawTag)

		option, err := parseConfigOption(tag)
		if err != nil && errors.Is(err, errKeyNotFound) {
			continue
		}
		if err != nil {
			logger.Warnw("failed to parse struct tag into ConfigOption", "err", err)
			continue
		}

		options = append(options, option)
	}

	return options
}
