package types

// Tag struct contaings information of a experiment run tag
type Tag struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
	PluginName  string `json:"pluginName"`
}
