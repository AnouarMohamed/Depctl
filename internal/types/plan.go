package types

// ProjectPlan stores metadata about the project being planned for deployment.
type ProjectPlan struct {
	Name string `json:"name"`
	Root string `json:"root"`
}

// RuntimePlan represents a simplified view of the runtime to be configured.
type RuntimePlan struct {
	Name       string  `json:"name"`
	Framework  string  `json:"framework"`
	Confidence float64 `json:"confidence"`
}

// Service represents a runtime component (e.g. app container, db container) in the deployment configuration.
type Service struct {
	Name         string `json:"name"`
	Type         string `json:"type"` // e.g. "app", "database"
	Build        string `json:"build,omitempty"`
	Image        string `json:"image,omitempty"`
	InternalPort int    `json:"internalPort,omitempty"`
	Public       bool   `json:"public"`
	Volume       string `json:"volume,omitempty"`
}

// EnvPlan models the environment requirements derived from scanning.
type EnvPlan struct {
	Required  []string `json:"required"`
	Sensitive []string `json:"sensitive"`
}

// Action represents an deployment engine step (e.g. network creation, compose invocation).
type Action struct {
	Type string `json:"type"` // e.g. "create_network", "compose_up"
	Name string `json:"name,omitempty"`
	File string `json:"file,omitempty"`
}

// Plan is the complete execution blueprint used by depctl to orchestrate deployment.
type Plan struct {
	Version        string            `json:"version"`
	Project        ProjectPlan       `json:"project"`
	Preset         string            `json:"preset"`
	Domain         string            `json:"domain"`
	PublicService  string            `json:"publicService"`
	Runtime        RuntimePlan       `json:"runtime"`
	Services       []Service         `json:"services"`
	Env            EnvPlan           `json:"env"`
	GeneratedFiles []string          `json:"generatedFiles"`
	Actions        []Action          `json:"actions"`
	Warnings       []string          `json:"warnings,omitempty"`
	ManualSteps    []string          `json:"manualSteps,omitempty"`
	FileHashes     map[string]string `json:"fileHashes,omitempty"`
}
