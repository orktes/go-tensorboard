package types

/*
{"data_location": "experiment AdYd1TgeTlaLWXx6I8JUbA", "window_title": "", "experiment_name": "My latest experiment", "experiment_description": "Simple comparison of several hyperparameters", "creation_time": 1572312164}
*/

// Environment struct containing current environment information
type Environment struct {
	WindowTime            string `json:"window_time"`
	DataLocation          string `json:"data_location"`
	ExperimentName        string `json:"experiment_name"`
	ExperimentDescription string `json:"experiment_description"`
	CreationTime          int    `json:"creation_time"`
}
