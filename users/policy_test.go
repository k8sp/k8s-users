package users

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultServiceAcccount(t *testing.T) {
	namespace := getDefaultServiceAcccount("test")
	assert.True(t, strings.EqualFold(namespace, "system:serviceaccount:test:default"))
}
