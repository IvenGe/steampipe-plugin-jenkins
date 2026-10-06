package jenkins

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func TestSplitFullName(t *testing.T) {
	tests := []struct {
		fullName string
		name     string
		parents  []string
	}{
		{fullName: "job", name: "job", parents: []string{}},
		{fullName: "folder/job", name: "job", parents: []string{"folder"}},
		{fullName: "a/b/c", name: "c", parents: []string{"a", "b"}},
		{fullName: "/folder/job/", name: "job", parents: []string{"folder"}},
		{fullName: "", name: "", parents: nil},
		{fullName: "/", name: "", parents: nil},
	}

	for _, tt := range tests {
		t.Run(tt.fullName, func(t *testing.T) {
			name, parents := splitFullName(tt.fullName)
			assert.Equal(t, tt.name, name)
			if tt.parents == nil {
				assert.Nil(t, parents)
			} else {
				assert.Equal(t, tt.parents, parents)
			}
		})
	}
}

func TestIsNotFoundErr(t *testing.T) {
	assert.False(t, isNotFoundErr(nil))
	assert.True(t, isNotFoundErr(errors.New("404")))
	assert.True(t, isNotFoundErr(errors.New("Build not found")))
	assert.True(t, isNotFoundErr(errors.New("No node found")))
	assert.False(t, isNotFoundErr(errors.New("connection refused")))
	// Legacy incorrect pattern should not be the only matcher anymore
	assert.False(t, isNotFoundErr(errors.New("Not found")))
}

func TestIsNotFoundErrorPredicate(t *testing.T) {
	pred := isNotFoundError(nil)
	assert.True(t, pred(context.Background(), nil, nil, errors.New("404")))
	assert.False(t, pred(context.Background(), nil, nil, errors.New("boom")))
	assert.False(t, pred(context.Background(), nil, nil, nil))
}

func TestInt64FromMap(t *testing.T) {
	m := map[string]interface{}{
		"Number": float64(42),
		"Alt":    int64(7),
		"Int":    3,
	}
	n, ok := int64FromMap(m, "Number")
	require.True(t, ok)
	assert.Equal(t, int64(42), n)

	n, ok = int64FromMap(m, "Alt")
	require.True(t, ok)
	assert.Equal(t, int64(7), n)

	n, ok = int64FromMap(m, "Int")
	require.True(t, ok)
	assert.Equal(t, int64(3), n)

	_, ok = int64FromMap(m, "missing")
	assert.False(t, ok)
	_, ok = int64FromMap(nil, "Number")
	assert.False(t, ok)
}

func TestGetConfigHandlesPointerAndValue(t *testing.T) {
	url := "https://jenkins.example"
	user := "admin"
	pass := "secret"

	valueCfg := jenkinsConfig{ServerURL: &url, Username: &user, Password: &pass}
	got := GetConfig(&plugin.Connection{Config: valueCfg})
	require.NotNil(t, got.ServerURL)
	assert.Equal(t, url, *got.ServerURL)

	ptrCfg := &jenkinsConfig{ServerURL: &url, Username: &user, Password: &pass}
	got = GetConfig(&plugin.Connection{Config: ptrCfg})
	require.NotNil(t, got.Username)
	assert.Equal(t, user, *got.Username)

	got = GetConfig(nil)
	assert.Nil(t, got.ServerURL)
}
