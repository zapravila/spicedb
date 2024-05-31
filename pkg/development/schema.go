package development

import (
	"errors"

	devinterface "github.com/zapravila/spicedb/pkg/proto/developer/v1"
	"github.com/zapravila/spicedb/pkg/schemadsl/compiler"
	"github.com/zapravila/spicedb/pkg/schemadsl/input"
)

// CompileSchema compiles a schema into its caveat and namespace definition(s), returning a developer
// error if the schema could not be compiled. The non-developer error is returned only if an
// internal errors occurred.
func CompileSchema(schema string) (*compiler.CompiledSchema, *devinterface.DeveloperError, error) {
	compiled, err := compiler.Compile(compiler.InputSchema{
		Source:       input.Source("schema"),
		SchemaString: schema,
	}, compiler.AllowUnprefixedObjectType())

	var contextError compiler.ErrorWithContext
	if errors.As(err, &contextError) {
		line, col, lerr := contextError.SourceRange.Start().LineAndColumn()
		if lerr != nil {
			return nil, nil, lerr
		}

		return nil, &devinterface.DeveloperError{
			Message: contextError.BaseCompilerError.BaseMessage,
			Kind:    devinterface.DeveloperError_SCHEMA_ISSUE,
			Source:  devinterface.DeveloperError_SCHEMA,
			Line:    uint32(line) + 1, // 0-indexed in parser.
			Column:  uint32(col) + 1,  // 0-indexed in parser.
			Context: contextError.ErrorSourceCode,
		}, nil
	}

	if err != nil {
		return nil, nil, err
	}

	return compiled, nil, nil
}
