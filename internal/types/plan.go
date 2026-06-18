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

// TargetPlan describes where the deployment should run.
type TargetPlan struct {
	Kind      string `json:"kind"`                // vps, vercel, fly
	Preset    string `json:"preset,omitempty"`    // compose-traefik, compose-nginx
	Root      string `json:"root,omitempty"`      // project root used by provider commands
	OutputDir string `json:"outputDir,omitempty"` // audit directory, usually .deploy
	AppName   string `json:"appName,omitempty"`
	Region    string `json:"region,omitempty"`
	EnvFile   string `json:"envFile,omitempty"`
}

// Artifact models a generated file and whether it belongs in .deploy or the repo root.
type Artifact struct {
	Path     string `json:"path"`
	Kind     string `json:"kind,omitempty"`     // compose, dockerfile, provider-config, report, script
	Scope    string `json:"scope,omitempty"`    // deploy, root
	Template string `json:"template,omitempty"` // internal template identifier
	Mode     string `json:"mode,omitempty"`     // file mode such as 0644 or 0755
}

// Check describes a preflight or validation check that should pass before apply.
type Check struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

// CredentialRequirement describes external auth expected by a provider.
type CredentialRequirement struct {
	Name        string `json:"name"`
	EnvVar      string `json:"envVar,omitempty"`
	Command     string `json:"command,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
}

// SecretImport describes how provider apply should import secret values.
// It stores only filenames and key names, never secret values.
type SecretImport struct {
	SourceFile string   `json:"sourceFile"`
	Keys       []string `json:"keys"`
	Mode       string   `json:"mode"` // import-env-file, stdin-per-key
}

// RollbackPlan stores provider-specific rollback metadata.
type RollbackPlan struct {
	Strategy     string `json:"strategy"`
	StateFile    string `json:"stateFile,omitempty"`
	LastImage    string `json:"lastImage,omitempty"`
	LastDeployID string `json:"lastDeployId,omitempty"`
}

// Action represents an deployment engine step (e.g. network creation, compose invocation).
type Action struct {
	Type    string   `json:"type"` // e.g. "create_network", "compose_up"
	Name    string   `json:"name,omitempty"`
	File    string   `json:"file,omitempty"`
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
	Cwd     string   `json:"cwd,omitempty"`
}

// Plan is the complete execution blueprint used by depctl to orchestrate deployment.
type Plan struct {
	Version        string                  `json:"version"`
	Target         TargetPlan              `json:"target"`
	Project        ProjectPlan             `json:"project"`
	Preset         string                  `json:"preset"`
	Domain         string                  `json:"domain"`
	PublicService  string                  `json:"publicService"`
	Runtime        RuntimePlan             `json:"runtime"`
	Build          BuildDetection          `json:"build"`
	Network        NetworkDetection        `json:"network"`
	Services       []Service               `json:"services"`
	Env            EnvPlan                 `json:"env"`
	GeneratedFiles []string                `json:"generatedFiles"`
	Artifacts      []Artifact              `json:"artifacts,omitempty"`
	Checks         []Check                 `json:"checks,omitempty"`
	Actions        []Action                `json:"actions"`
	Credentials    []CredentialRequirement `json:"credentials,omitempty"`
	SecretImports  []SecretImport          `json:"secretImports,omitempty"`
	Rollback       RollbackPlan            `json:"rollback,omitempty"`
	Warnings       []string                `json:"warnings,omitempty"`
	ManualSteps    []string                `json:"manualSteps,omitempty"`
	FileHashes     map[string]string       `json:"fileHashes,omitempty"`
}
