package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/orktes/go-tensorboard/handler"
	"github.com/orktes/go-tensorboard/plugins"
	"github.com/orktes/go-tensorboard/plugins/scalars"
	"github.com/orktes/go-tensorboard/types"
)

type mockDataloader struct {
}

func (mdl *mockDataloader) ListRuns(ctx context.Context) ([]string, error) {
	return []string{"test1", "test2"}, nil
}

func (mdl *mockDataloader) GetEnvironment(ctx context.Context) (types.Environment, error) {
	return types.Environment{
		WindowTime:            "Experiment",
		DataLocation:          "Experiment",
		ExperimentName:        "Experiment",
		ExperimentDescription: "Experiment",
		CreationTime:          int(time.Now().Unix()),
	}, nil
}

func (mdl *mockDataloader) GetPluginTags(ctx context.Context, pluginName string) (types.PluginRunTags, error) {
	return types.PluginRunTags{
		"test1": types.PluginTags{
			"accuracy/accuracy": map[string]interface{}{"displayName": "accuracy/accuracy", "description": "Model accuracy"},
		},
		"test2": types.PluginTags{
			"accuracy/accuracy": map[string]interface{}{"displayName": "accuracy/accuracy", "description": "Model accuracy"},
		},
	}, nil
}

func (mdl *mockDataloader) GetPluginData(ctx context.Context, pluginName string, resource string, query types.PluginQuery) (interface{}, error) {
	run := query["run"]
	multiplier := 1

	if run == "test2" {
		multiplier = 2
	}

	return []scalars.ScalarValue{
		scalars.ScalarValue{
			WallTime: time.Now(),
			Step:     1,
			Value:    float64(multiplier * 1),
		},
		scalars.ScalarValue{
			WallTime: time.Now().Add(time.Second),
			Step:     2,
			Value:    float64(multiplier * 3),
		},
		scalars.ScalarValue{
			WallTime: time.Now().Add(2 * time.Second),
			Step:     3,
			Value:    float64(multiplier * 4),
		},
	}, nil
}

func main() {
	h := handler.New(&mockDataloader{}, plugins.DefaultRegistry)
	http.Handle("/", h)

	fmt.Println("Starting in :8080. Visit http://localhost:8080/")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
