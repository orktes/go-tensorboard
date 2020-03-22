package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"

	"github.com/orktes/go-tensorboard/types"
	"github.com/orktes/go-tensorboard/ui"
)

// DataLoader describes the interface for dataloader
type DataLoader interface {
	ListRuns(ctx context.Context, experimentID string) ([]string, error)
	GetEnvironment(ctx context.Context, experimentID string) (types.Environment, error)
	GetPluginTags(ctx context.Context, experimentID string, pluginName string) (types.PluginRunTags, error)
	GetPluginData(ctx context.Context, experimentID string, pluginName string, resource string, query types.PluginQuery) (interface{}, error)
}

// PluginLoader describes the interface for the plugin loader
type PluginLoader interface {
	PluginEntry(ctx context.Context, name string) ([]byte, error)
	ListPlugins(ctx context.Context) (map[string]types.PluginConfig, error)
}

var _ http.Handler = &Handler{}

// Handler implements the apis for a Tensorboard backend
type Handler struct {
	dataLoader   DataLoader
	pluginLoader PluginLoader
	fileServer   http.Handler
	*mux.Router
}

// New returns a new handler for given DataLoader and PluginLoader
func New(dl DataLoader, pl PluginLoader) http.Handler {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    ui.Asset,
		AssetDir: ui.AssetDir,
	})

	h := &Handler{
		dataLoader:   dl,
		pluginLoader: pl,
		fileServer:   fileServer,
		Router:       mux.NewRouter(),
	}

	h.initRoutes()

	return h
}

func (h *Handler) handleLogDir(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"logdir": r.Referer()})
}

func (h *Handler) handlePluginListing(w http.ResponseWriter, r *http.Request) error {
	res, err := h.pluginLoader.ListPlugins(r.Context())
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(res)
}

func (h *Handler) handleEnvironment(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	env, err := h.dataLoader.GetEnvironment(r.Context(), vars["experimentID"])
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(env)
}

func (h *Handler) handleRuns(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	runs, err := h.dataLoader.ListRuns(r.Context(), vars["experimentID"])
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(runs)
}

func (h *Handler) handlePluginTags(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	tags, err := h.dataLoader.GetPluginTags(r.Context(), vars["experimentID"], vars["pluginName"])
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(tags)
}

func (h *Handler) handlePluginData(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	pq := types.PluginQuery{}

	q := r.URL.Query()

	for key, val := range q {
		pq[key] = val[0]
	}

	data, err := h.dataLoader.GetPluginData(
		r.Context(), vars["experimentID"], vars["pluginName"], vars["data"], pq)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(data)
}

func (h *Handler) handleExperiments(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("[]"))
	return nil
}

func (h *Handler) handleUI(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")

	path := strings.Join(pathParts[2:], "/")
	r2 := new(http.Request)
	*r2 = *r
	r2.URL = new(url.URL)
	*r2.URL = *r.URL
	r2.URL.Path = path

	h.fileServer.ServeHTTP(w, r2)
}

func (h *Handler) e(hf func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := hf(w, r)
		if err != nil {
			// TODO handle err
			fmt.Printf("%+v\n", err)
		}
	}
}

func (h *Handler) initRoutes() {
	h.HandleFunc("/{experimentID}/data/logdir", h.handleLogDir)
	h.HandleFunc("/{experimentID}/data/plugins_listing", h.e(h.handlePluginListing))
	h.HandleFunc("/{experimentID}/data/environment", h.e(h.handleEnvironment))
	h.HandleFunc("/{experimentID}/data/runs", h.e(h.handleRuns))
	h.HandleFunc("/{experimentID}/data/experiments", h.e(h.handleExperiments))
	h.HandleFunc("/{experimentID}/data/plugin/{pluginName}/tags", h.e(h.handlePluginTags))
	h.HandleFunc("/{experimentID}/data/plugin/{pluginName}/{data}", h.e(h.handlePluginData))
	h.PathPrefix("/{experimentID}/").HandlerFunc(h.handleUI)
}
