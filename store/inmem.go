package store

import (
	"context"
	"errors"
	"io"
	"math"
	"time"

	"github.com/Applifier/go-tensorflow/types/tensorflow/core/framework"
	"github.com/orktes/go-tensorboard/events"
	"github.com/orktes/go-tensorboard/plugins/scalars"
	"github.com/orktes/go-tensorboard/types"
)

// InMemStore in memory summary event storage
type InMemStore struct {
	// scalars map from tag to scalar values
	scalars map[string][]scalars.Value
}

// NewInMemStore returns a new inmem storage instance
func NewInMemStore() *InMemStore {
	return &InMemStore{
		scalars: map[string][]scalars.Value{},
	}
}

func (inmem *InMemStore) processScalarValue(tag string, wallTime float64, step int, v *framework.Summary_Value) {
	scalarValue := 0.0
	switch val := v.Value.(type) {
	case *framework.Summary_Value_SimpleValue:
		scalarValue = float64(val.SimpleValue)
	case *framework.Summary_Value_Tensor:
		tensorValue, err := events.TensorContentToGoType(val.Tensor)
		if err != nil {
			panic(err)
		}
		switch tVal := tensorValue.(type) {
		case float32:
			scalarValue = float64(tVal)
		case float64:
			scalarValue = float64(tVal)
		case int:
			scalarValue = float64(tVal)
		case int32:
			scalarValue = float64(tVal)
		case int64:
			scalarValue = float64(tVal)
		}
	}

	values := inmem.scalars[tag]
	i, f := math.Modf(wallTime)
	values = append(values, scalars.Value{
		WallTime: time.Unix(int64(i), int64(f*1e9)),
		Step:     step,
		Value:    scalarValue,
	})
	inmem.scalars[tag] = values
}

// GetPluginTags returns tags for given plugin
func (inmem *InMemStore) GetPluginTags(ctx context.Context, pluginName string) (types.PluginTags, error) {
	tags := types.PluginTags{}

	switch pluginName {
	case "scalars":
		for key := range inmem.scalars {
			tags[key] = map[string]interface{}{"displayName": key, "description": ""}
		}
	}

	return tags, nil
}

// GetPluginData returns plugin data
func (inmem *InMemStore) GetPluginData(ctx context.Context, pluginName string, resource string, query types.PluginQuery) (interface{}, error) {
	switch pluginName {
	case "scalars":
		if resource == "scalars" {
			vals := inmem.scalars[query["tag"]]
			if len(vals) == 0 {
				vals = []scalars.Value{}
			}
			return vals, nil
		}
	}

	return nil, errors.New("unknown plugin")
}

// Populate polulates inmemory storage from given events.Reader
func (inmem *InMemStore) Populate(r *events.Reader) error {
	for {
		ev, err := r.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		summary := ev.GetSummary()
		if summary == nil {
			continue
		}

		for _, v := range summary.Value {

			if v.Metadata != nil && v.Metadata.PluginData != nil {
				pluginName := v.Metadata.PluginData.PluginName
				switch pluginName {
				case "scalars":
					inmem.processScalarValue(v.Tag, ev.WallTime, int(ev.Step), v)
				}
			}
		}
	}
}
