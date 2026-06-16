package planner

import (
	"fmt"
	"strings"

	"depctl/internal/types"
)

// GeneratePlanReport creates a markdown report detailing the compiled deployment plan.
func GeneratePlanReport(plan *types.Plan) string {
	var sb strings.Builder

	sb.WriteString("# Deployment Plan Report\n\n")
	sb.WriteString(fmt.Sprintf("- **Preset:** %s\n", plan.Preset))
	sb.WriteString(fmt.Sprintf("- **Target Domain:** %s\n", plan.Domain))
	sb.WriteString(fmt.Sprintf("- **Public Service:** %s\n\n", plan.PublicService))

	sb.WriteString("## Services Setup\n")
	for _, svc := range plan.Services {
		sb.WriteString(fmt.Sprintf("### Service: %s\n", svc.Name))
		sb.WriteString(fmt.Sprintf("- **Type:** %s\n", svc.Type))
		if svc.Image != "" {
			sb.WriteString(fmt.Sprintf("- **Image:** %s\n", svc.Image))
		}
		if svc.Build != "" {
			sb.WriteString(fmt.Sprintf("- **Build Context:** %s\n", svc.Build))
		}
		if svc.InternalPort != 0 {
			sb.WriteString(fmt.Sprintf("- **Internal Port:** %d\n", svc.InternalPort))
		}
		if svc.Volume != "" {
			sb.WriteString(fmt.Sprintf("- **Volume:** %s\n", svc.Volume))
		}
		sb.WriteString(fmt.Sprintf("- **Publicly Exposed:** %v\n\n", svc.Public))
	}

	sb.WriteString("## Environment Requirements\n")
	if len(plan.Env.Required) > 0 {
		for _, k := range plan.Env.Required {
			isSens := false
			for _, s := range plan.Env.Sensitive {
				if s == k {
					isSens = true
					break
				}
			}
			if isSens {
				sb.WriteString(fmt.Sprintf("- %s *(Sensitive Key)*\n", k))
			} else {
				sb.WriteString(fmt.Sprintf("- %s\n", k))
			}
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("- None\n\n")
	}

	sb.WriteString("## Files to Generate\n")
	for _, f := range plan.GeneratedFiles {
		sb.WriteString(fmt.Sprintf("- %s\n", f))
	}
	sb.WriteString("\n")

	sb.WriteString("## Actions to Execute on Apply\n")
	for i, act := range plan.Actions {
		sb.WriteString(fmt.Sprintf("%d. **%s**", i+1, act.Type))
		if act.Name != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", act.Name))
		}
		if act.File != "" {
			sb.WriteString(fmt.Sprintf(" using %s", act.File))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	if len(plan.Warnings) > 0 {
		sb.WriteString("## Plan Warnings\n")
		for _, w := range plan.Warnings {
			sb.WriteString(fmt.Sprintf("- %s\n", w))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Required Manual Steps\n")
	for i, step := range plan.ManualSteps {
		sb.WriteString(fmt.Sprintf("%d. [ ] %s\n", i+1, step))
	}
	sb.WriteString("\n")

	sb.WriteString("## Next Steps\n")
	sb.WriteString("Run the writer command to compile the deployment files:\n")
	sb.WriteString("```bash\ndepctl write\n```\n")

	return sb.String()
}
