package types

// Evidence represents a single piece of evidence supporting a detection.
type Evidence struct {
	File   string `json:"file"`
	Reason string `json:"reason"`
}

// ProjectDetection contains metadata about the analyzed project.
type ProjectDetection struct {
	Name    string `json:"name"`
	Root    string `json:"root"`
	GitRoot bool   `json:"gitRoot"`
}

// RuntimeDetection holds information about the detected runtime environment (e.g. Node, Python).
type RuntimeDetection struct {
	Name       string     `json:"name"`
	Framework  string     `json:"framework"`
	Confidence float64    `json:"confidence"`
	Evidence   []Evidence `json:"evidence,omitempty"`
}

// BuildDetection contains commands and tools used for building/running the application.
type BuildDetection struct {
	PackageManager string `json:"packageManager"`
	BuildCommand   string `json:"buildCommand"`
	StartCommand   string `json:"startCommand"`
}

// NetworkDetection holds the internal port configuration and its confidence metrics.
type NetworkDetection struct {
	InternalPort int        `json:"internalPort"`
	Confidence   float64    `json:"confidence"`
	Evidence     []Evidence `json:"evidence,omitempty"`
}

// EnvDetection records environment variable keys, sensitive keys, and env file presence.
type EnvDetection struct {
	Keys          []string `json:"keys"`
	Sensitive     []string `json:"sensitive"`
	HasEnvExample bool     `json:"hasEnvExample"`
}

// Dependency details the likelihood and confidence level of standard service dependencies.
type Dependency struct {
	Likely     bool    `json:"likely"`
	Confidence float64 `json:"confidence"`
}

// ContainerizationDetection flags existing deployment configurations like Dockerfiles or Compose files.
type ContainerizationDetection struct {
	DockerfilePresent bool `json:"dockerfilePresent"`
	ComposePresent    bool `json:"composePresent"`
}

// CIDetection checks for existing CI pipeline structures.
type CIDetection struct {
	GitHubActions bool `json:"githubActions"`
	GitLab        bool `json:"gitlab"`
	GitEA         bool `json:"gitea"`
}

// Detection is the root data model containing all results parsed from a project scan.
type Detection struct {
	Version          string                    `json:"version"`
	Project          ProjectDetection          `json:"project"`
	Runtime          RuntimeDetection          `json:"runtime"`
	Build            BuildDetection            `json:"build"`
	Network          NetworkDetection          `json:"network"`
	Env              EnvDetection              `json:"env"`
	Dependencies     map[string]Dependency     `json:"dependencies"`
	Containerization ContainerizationDetection `json:"containerization"`
	CI               CIDetection               `json:"ci"`
	Warnings         []string                  `json:"warnings,omitempty"`
}
