package api

type LoadTest struct {
	Target  Target  `json:"target"`
	Context Context `json:"context"`
	Steps   []Step  `json:"steps"`
}

type Target struct {
	Protocol string            `json:"protocol"`
	BaseURL  string            `json:"base_url"`
	Endpoint DynamicString     `json:"endpoint"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
}

type DynamicString struct {
	Format  string  `json:"format"`
	Context Context `json:"context"`
}

type Context [][]interface{}

type Step struct {
	RPMFrom     int `json:"rpm_from"`
	RPMTo       int `json:"rpm_to"`
	DurationSec int `json:"duration_seconds"`
}

type Result struct {
	StepResults     []StepResult `json:"step_results"`
	DurationSeconds float64      `json:"duration_seconds"`
}

type StepResult struct {
	CallCount       int          `json:"call_count"`
	Status          StatusResult `json:"status"`
	DurationSeconds float64      `json:"duration_seconds"`
}

type StatusResult map[int]int
