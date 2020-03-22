package types

// ExperimentRun struct contains information for an experiment run
type ExperimentRun struct {
	ID         int    `json:"id"`
	Names      string `json:"names"`
	StartTimeS int    `json:"startTime"`
	Tags       []Tag  `json:"tags"`
}
