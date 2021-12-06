package setup

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"

	"edholm.dev/go-logging"
)

func ParseAst(ctx context.Context, filenames []string) []*ast.File {
	logger := logging.FromContext(ctx)

	files := make([]*ast.File, 0, len(filenames))
	set := token.NewFileSet()
	for _, file := range filenames {
		parsedFile, err := parser.ParseFile(set, file, nil, 0)
		if err != nil {
			logger.Info("failed to parse file", "file", file, "err", err)
			continue
		}

		files = append(files, parsedFile)
	}

	return files
}
