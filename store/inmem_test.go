package store

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/orktes/go-tensorboard/events"
	"github.com/stretchr/testify/require"
)

func TestInMemStore(t *testing.T) {
	testFile := path.Join(getTestDataDir(), "tfevents/events.out.tfevents.1585499054.Jaakkos-Mac-mini.local.11684.5.v2")
	f, err := os.Open(testFile)
	require.NoError(t, err)
	defer f.Close()
	r := events.NewReader(f)
	store := NewInMemStore()

	t.Run("should consume reader successfully", func(t *testing.T) {
		err := store.Populate(r)
		require.NoError(t, err)
	})

	t.Run("should have 100 scalars", func(t *testing.T) {
		scalars := store.scalars["my_metric"]
		require.Len(t, scalars, 100)

		require.Equal(t, float64(0), scalars[0].Value)
		require.Equal(t, float64(49.5), scalars[99].Value)
	})
}

func getTestPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filename
}

func getTestDataDir() string {
	return path.Join(getTestPath(), "../../testdata")
}
