package plugins

import (
	"context"
	"errors"

	"github.com/orktes/go-tensorboard/types"
)

// DefaultRegistry contains default plugin registry
var DefaultRegistry = &Registry{}

// Registry plugin registry
type Registry struct{}

// PluginEntry returns pluging entry as html
func (r *Registry) PluginEntry(ctx context.Context, name string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

// ListPlugins returns plugin configuration for installed plugins
func (r *Registry) ListPlugins(ctx context.Context) (map[string]types.PluginConfig, error) {
	return map[string]types.PluginConfig{
		"scalars": types.PluginConfig{
			Enabled:   true,
			RemoveDom: false,
			TabName:   "scalars",
			LoadingMechanism: types.PluginLoadingMechanism{
				Type:        "CUSTOM_ELEMENT",
				ElementName: "tf-scalar-dashboard",
			},
		},
	}, nil
}
