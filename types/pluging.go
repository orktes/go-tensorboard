package types

// PluginData data returned based on plugin
type PluginData interface{}

// PluginQuery map contains call arguments for plugin calls
type PluginQuery map[string]string

// PluginTags map containing plugin tag information
type PluginTags map[string]map[string]interface{}

// PluginRunTags map of to experiment runs runs to plugi tags
type PluginRunTags map[string]PluginTags

// PluginLoadingMechanism describes how plugin should be loaded
type PluginLoadingMechanism struct {
	Type        string `json:"type"`
	ElementName string `json:"element_name,omitempty"`
}

// PluginConfig struct contains information about a configured plugin
type PluginConfig struct {
	DisableReload    bool                   `json:"disable_reload"`
	Enabled          bool                   `json:"enabled"`
	RemoveDom        bool                   `json:"remove_dom"`
	TabName          string                 `json:"tab_name"`
	LoadingMechanism PluginLoadingMechanism `json:"loading_mechanism"`
}
