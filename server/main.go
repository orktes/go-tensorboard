package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/orktes/go-tensorboard/events"
	"github.com/orktes/go-tensorboard/handler"
	"github.com/orktes/go-tensorboard/plugins"
	"github.com/orktes/go-tensorboard/store"
	"github.com/orktes/go-tensorboard/types"
)

type mockDataloader struct {
	store *store.InMemStore
}

func (mdl *mockDataloader) ListRuns(ctx context.Context) ([]string, error) {
	return []string{"default"}, nil
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
	tags, err := mdl.store.GetPluginTags(ctx, pluginName)
	if err != nil {
		return nil, err
	}
	return types.PluginRunTags{
		"default": tags,
	}, nil
}

func (mdl *mockDataloader) GetPluginData(ctx context.Context, pluginName string, resource string, query types.PluginQuery) (interface{}, error) {
	return mdl.store.GetPluginData(ctx, pluginName, resource, query)
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "" {
		fmt.Printf("give command the logdir as the first argument (%s path_to_logir)", os.Args[0])
		os.Exit(1)
	}

	logdir := os.Args[1]

	store := store.NewInMemStore()

	dir := path.Join(logdir, "*tfevents*")

	m, err := filepath.Glob(dir)
	if err != nil {
		panic(err)
	}

	for _, p := range m {
		f, err := os.Open(p)
		if err != nil {
			panic(err)
		}

		r := events.NewReader(f)

		if err := store.Populate(r); err != nil {
			panic(err)
		}

		f.Close()
	}

	h := handler.New(&mockDataloader{store: store}, plugins.DefaultRegistry)
	http.Handle("/", h)

	fmt.Println("Starting in :8080. Visit http://localhost:8080/")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
