package ports

import "io"

type DemotivatorGenerator interface {
	Generate(resultWriter io.Writer, data any) error
}
