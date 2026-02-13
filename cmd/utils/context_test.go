//go:build !integration

package utils

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	assert.Equal(t, viper.GetViper(), ViperFromContext(t.Context()))

	newV := viper.New()
	ctx := WithViper(t.Context(), newV)

	newV.Set("test", "test")
	assert.Equal(t, newV, ViperFromContext(ctx))
}
