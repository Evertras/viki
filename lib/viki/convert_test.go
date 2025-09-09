package viki_test

import (
	"testing"

	"github.com/evertras/viki/lib/viki"
	"github.com/stretchr/testify/assert"
)

func TestNewConverter(t *testing.T) {
	converter := viki.NewConverter(viki.ConverterOptions{})
	assert.NotNil(t, converter)
}
