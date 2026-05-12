package parse

import (
	"github.com/mtraver/matrixfont/log"
)

type parseOptions struct {
	verbosity int
}

func defaultOptions() *parseOptions {
	return &parseOptions{
		verbosity: log.LevelWarn,
	}
}

type Opt func(opts *parseOptions) error

func WithLogVerbosity(verbosity int) Opt {
	return func(opts *parseOptions) error {
		opts.verbosity = verbosity
		return nil
	}
}
