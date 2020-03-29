package events

import (
	"io"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/Applifier/go-tensorflow/types/tensorflow/core/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	testFile := path.Join(getTestDataDir(), "tfevents/events.out.tfevents.1585499054.Jaakkos-Mac-mini.local.11684.5.v2")
	f, err := os.Open(testFile)
	require.NoError(t, err)
	defer f.Close()

	r := NewReader(f)

	events := make([]*util.Event, 0, 99)

	t.Run("first event should util.Event_FileVersion", func(t *testing.T) {
		ev, err := r.Next()
		require.NoError(t, err)

		fv := ev.GetFileVersion()
		require.Equal(t, "brain.Event:2", fv)
	})

	t.Run("should contain summary events", func(t *testing.T) {
		ev, err := r.Next()
		require.NoError(t, err)

		summary := ev.GetSummary()

		value := summary.Value[0]
		assert.Equal(t, "my_metric", value.Tag)
		assert.Equal(t, "scalars", value.Metadata.PluginData.PluginName)
		assert.Equal(t, int64(0), ev.Step)

		tensor := value.GetTensor()

		val, err := TensorContentToGoType(tensor)
		require.NoError(t, err)

		require.Equal(t, float32(0), val)

	})

	t.Run("should read rest of 99 events from reader", func(t *testing.T) {
		for {
			ev, err := r.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				require.NoError(t, err)
			}

			events = append(events, ev)
		}

		require.Len(t, events, 99, "should contain 99 events")
	})

}

func getTestPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filename
}

func getTestDataDir() string {
	return path.Join(getTestPath(), "../../testdata")
}
