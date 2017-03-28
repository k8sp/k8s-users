package users

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultServiceAcccount(t *testing.T) {
	username := getDefaultServiceAcccount("test")
	assert.True(t, strings.EqualFold(username, "system:serviceaccount:test:default"))
}
