package environ_test

import (
	"os"
	"strings"
	"testing"

	"code.olapie.com/environ"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvVar(t *testing.T) {
	key := "A1_B2_C3"
	val := uuid.New().String()
	err := os.Setenv(key, val)
	require.NoError(t, err)

	assert.Equal(t, val, environ.String(key, ""))
	valByLowercase := environ.String(strings.ToLower(key), "")
	assert.Equal(t, val, valByLowercase)
	valByDotKey := environ.String(strings.Replace(strings.ToLower(key), "_", ".", -1), "")
	assert.Equal(t, val, valByDotKey)
}

// func TestTOML(t *testing.T) {
// 	err := environ.LoadConfigFile("testdata/config.toml")
// 	require.NoError(t, err)

// 	require.Equal(t, "user1", environ.String("alpha.database.user", ""))
// 	require.Equal(t, "password1", environ.String("alpha.database.password", ""))

// 	require.Equal(t, "user2", environ.String("beta.database.user", ""))
// 	require.Equal(t, "password2", environ.String("beta.database.password", ""))

// 	gammaHosts := environ.StringSlice("gamma.hosts", nil)
// 	require.Equal(t, []string{"gamma.host1", "gamma.host2"}, gammaHosts)
// }

func TestJSON(t *testing.T) {
	err := environ.LoadConfigFile("testdata/config.json")
	require.NoError(t, err)

	require.Equal(t, "json_user1", environ.String("alpha.database.user", ""))
	require.Equal(t, "json_password1", environ.String("alpha.database.password", ""))

	require.Equal(t, "json_user2", environ.String("beta.database.user", ""))
	require.Equal(t, "json_password2", environ.String("beta.database.password", ""))

	gammaHosts := environ.StringSlice("gamma.hosts", nil)
	require.Equal(t, []string{"gamma.json_host1", "gamma.json_host2"}, gammaHosts)
}
