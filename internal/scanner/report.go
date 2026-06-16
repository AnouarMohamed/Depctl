package scanner

import (
	"fmt"
	"strings"

	"depctl/internal/types"
)

// ConfidenceBand maps a float64 confidence score to a string representation.
func ConfidenceBand(confidence float64) string {
	if confidence >= 0.90 {
		return "strong"
	}
	if confidence >= 0.70 {
		return "likely"
	}
	if confidence >= 0.50 {
		return "weak"
	}
	return "unknown"
}

// GenerateScanReport converts the types.Detection struct into a human-readable markdown report.
func GenerateScanReport(det *types.Detection) string {
	var sb strings.Builder

	sb.WriteString("# Scan Report\n\n")
	sb.WriteString(fmt.Sprintf("## Project Context\n"))
	sb.WriteString(fmt.Sprintf("- **Name:** %s\n", det.Project.Name))
	sb.WriteString(fmt.Sprintf("- **Path:** %s\n", det.Project.Root))
	sb.WriteString(fmt.Sprintf("- **Git Repository:** %v\n\n", det.Project.GitRoot))

	sb.WriteString("## Runtime & Framework\n")
	frameworkStr := det.Runtime.Framework
	if frameworkStr == "" {
		frameworkStr = "None"
	}
	sb.WriteString(fmt.Sprintf("- **Runtime:** %s\n", det.Runtime.Name))
	sb.WriteString(fmt.Sprintf("- **Framework:** %s\n", frameworkStr))
	sb.WriteString(fmt.Sprintf("- **Confidence:** %s (%.0f%%)\n\n", ConfidenceBand(det.Runtime.Confidence), det.Runtime.Confidence*100))

	if len(det.Runtime.Evidence) > 0 {
		sb.WriteString("### Evidence:\n")
		for _, ev := range det.Runtime.Evidence {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", ev.File, ev.Reason))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Build & Expose Specs\n")
	sb.WriteString(fmt.Sprintf("- **Package Manager:** %s\n", det.Build.PackageManager))
	sb.WriteString(fmt.Sprintf("- **Build Command:** %s\n", det.Build.BuildCommand))
	sb.WriteString(fmt.Sprintf("- **Start Command:** %s\n", det.Build.StartCommand))
	sb.WriteString(fmt.Sprintf("- **Internal Service Port:** %d (Confidence: %s)\n\n", det.Network.InternalPort, ConfidenceBand(det.Network.Confidence)))

	sb.WriteString("## Service Dependencies\n")
	for name, dep := range det.Dependencies {
		if dep.Likely {
			sb.WriteString(fmt.Sprintf("- **%s:** Detected (Confidence: %s)\n", name, ConfidenceBand(dep.Confidence)))
		} else {
			sb.WriteString(fmt.Sprintf("- **%s:** Not detected\n", name))
		}
	}
	sb.WriteString("\n")

	sb.WriteString("## Docker & CI Environments\n")
	sb.WriteString(fmt.Sprintf("- **Existing Dockerfile:** %v\n", det.Containerization.DockerfilePresent))
	sb.WriteString(fmt.Sprintf("- **Existing Docker Compose:** %v\n", det.Containerization.ComposePresent))
	sb.WriteString(fmt.Sprintf("- **CI Pipelines:** GitHub Actions (%v), GitLab CI (%v), Gitea (%v)\n\n", det.CI.GitHubActions, det.CI.GitLab, det.CI.GitEA))

	if len(det.Warnings) > 0 {
		sb.WriteString("## Warnings\n")
		for _, w := range det.Warnings {
			sb.WriteString(fmt.Sprintf("- [ ] %s\n", w))
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("## Warnings\n- None\n\n")
	}

	return sb.String()
}
