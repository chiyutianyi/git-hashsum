package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	revision = "v1.0.0"
	revisionDate = "2022-11-11"

	assert.Equal(t, "v1.0.0 (2022-11-11)", Version())
}
