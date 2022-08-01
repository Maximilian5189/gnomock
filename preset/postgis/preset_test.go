package postgis_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgis"
	"github.com/stretchr/testify/require"
)

func TestPreset(t *testing.T) {
	t.Parallel()

	// TODO: use versions (tags) this preset supports
	for _, version := range []string{"11-2.5-alpine"} {
		t.Run(version, testPreset(version))
	}
}

func testPreset(version string) func(t *testing.T) {
	return func(t *testing.T) {
		p := postgis.Preset(
			postgis.WithUser("gnomock", "gnomick"),
			postgis.WithDatabase("mydb"),
			postgis.WithQueriesFile("./testdata/queries.sql"),
			postgis.WithVersion(version),
		)
		container, err := gnomock.Start(p)

		defer func() { require.NoError(t, gnomock.Stop(container)) }()

		require.NoError(t, err)

		addr := container.DefaultAddress()
		require.NotEmpty(t, addr)
	}
}

func TestPreset_withDefaults(t *testing.T) {
	t.Parallel()

	p := postgis.Preset()
	container, err := gnomock.Start(p)
	require.NoError(t, err)

	defer func() { require.NoError(t, gnomock.Stop(container)) }()

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s  dbname=%s sslmode=disable",
		container.Host, container.DefaultPort(),
		"gnomock", "gnomick", "mydb",
	)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	require.NoError(t, db.Close())
}

func TestPreset_wrongQueriesFile(t *testing.T) {
	t.Parallel()

	p := postgis.Preset(
		postgis.WithQueriesFile("./invalid"),
	)
	c, err := gnomock.Start(p)
	require.Error(t, err)
	require.Contains(t, err.Error(), "can't read queries file")
	require.NoError(t, gnomock.Stop(c))
}
