package planner

import (
	"testing"

	"depctl/internal/types"
)

func TestPlanGeneration(t *testing.T) {
	t.Run("Node + Traefik (No DB)", func(t *testing.T) {
		det := &types.Detection{
			Project: types.ProjectDetection{Name: "my-node-app", Root: "/srv/my-node-app"},
			Runtime: types.RuntimeDetection{Name: "node", Framework: "nextjs", Confidence: 0.95},
			Network: types.NetworkDetection{InternalPort: 3000, Confidence: 0.8},
			Dependencies: map[string]types.Dependency{
				"postgres": {Likely: false},
				"redis":    {Likely: false},
			},
		}

		plan, err := Plan(det, "compose-traefik", "app.example.com", "github")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if plan.Preset != "compose-traefik" {
			t.Errorf("expected preset %q, got %q", "compose-traefik", plan.Preset)
		}
		if plan.Domain != "app.example.com" {
			t.Errorf("expected domain %q, got %q", "app.example.com", plan.Domain)
		}
		if len(plan.Services) != 1 {
			t.Errorf("expected 1 service (app only), got %d", len(plan.Services))
		}
		if plan.Services[0].Name != "web" || plan.Services[0].InternalPort != 3000 {
			t.Errorf("unexpected service configuration: %+v", plan.Services[0])
		}
	})

	t.Run("Laravel + Traefik (With Postgres & Redis)", func(t *testing.T) {
		det := &types.Detection{
			Project: types.ProjectDetection{Name: "laravel-app", Root: "/srv/laravel-app"},
			Runtime: types.RuntimeDetection{Name: "laravel", Framework: "laravel", Confidence: 0.95},
			Network: types.NetworkDetection{InternalPort: 80, Confidence: 0.75},
			Dependencies: map[string]types.Dependency{
				"postgres": {Likely: true, Confidence: 0.8},
				"redis":    {Likely: true, Confidence: 0.8},
			},
		}

		plan, err := Plan(det, "compose-traefik", "laravel.example.com", "github")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Services should include web, postgres, and redis
		if len(plan.Services) != 3 {
			t.Errorf("expected 3 services, got %d", len(plan.Services))
		}

		hasPostgres := false
		hasRedis := false
		for _, svc := range plan.Services {
			if svc.Name == "postgres" && svc.Type == "database" {
				hasPostgres = true
			}
			if svc.Name == "redis" && svc.Type == "database" {
				hasRedis = true
			}
		}

		if !hasPostgres {
			t.Error("expected postgres service in plan")
		}
		if !hasRedis {
			t.Error("expected redis service in plan")
		}
	})

	t.Run("Error Empty Domain", func(t *testing.T) {
		det := &types.Detection{
			Project: types.ProjectDetection{Name: "empty-dom", Root: "/srv/empty-dom"},
		}
		_, err := Plan(det, "compose-traefik", "", "github")
		if err == nil {
			t.Error("expected error due to empty domain, got nil")
		}
	})
}
