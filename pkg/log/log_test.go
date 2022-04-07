package log

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WithName(t *testing.T) {
	defer Flush() // used for record logger printer

	logger := WithName("test")
	logger.Infow("Hello world!", "foo", "bar") // structed logger
}

func Test_WithValues(t *testing.T) {
	defer Flush() // used for record logger printer

	logger := WithValues("key", "value") // used for record context
	logger.Info("Hello world!")
	logger.Info("Hello world!")
}

func Test_V(t *testing.T) {
	defer Flush() // used for record logger printer

	V(0).Infow("Hello world!", "key", "value")
	V(1).Infow("Hello world!", "key", "value")
}

func Test_Option(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ExitOnError)
	opt := NewOptions()
	opt.AddFlags(fs)

	args := []string{"--log.level=debug"}
	err := fs.Parse(args)
	assert.Nil(t, err)

	assert.Equal(t, "debug", opt.Level)
}