package modelrouter

import (
	"context"
	"errors"
)

const (
	ProviderLocalOpenAICompatible = "local-openai-compatible"
	ProviderOllama                = "ollama"
	ProviderLMStudio              = "lmstudio"
	ProviderLlamaCPP              = "llama.cpp"
	ProviderSurrealML             = "surrealml"
	ProviderCloudOpenAICompatible = "cloud-openai-compatible"
)

const (
	TaskChat            = "chat"
	TaskExplain         = "explain"
	TaskPlan            = "plan"
	TaskGenerateYAML    = "generate_yaml"
	TaskGenerateRego    = "generate_rego"
	TaskGenerateSurreal = "generate_surrealql"
	TaskClassify        = "classify"
	TaskScoreRisk       = "score_risk"
	TaskAnomalyDetect   = "anomaly_detect"
)

// Request is the normalized input passed to the model router.
type Request struct {
	Task          string         `json:"task"`
	Prompt        string         `json:"prompt,omitempty"`
	SystemPrompt  string         `json:"systemPrompt,omitempty"`
	Context       map[string]any `json:"context,omitempty"`
	PreferredModel string        `json:"preferredModel,omitempty"`
	AllowCloud    bool           `json:"allowCloud"`
}

// Response is the normalized model-router response.
type Response struct {
	Provider string         `json:"provider"`
	Model    string         `json:"model"`
	Text     string         `json:"text,omitempty"`
	Signals  map[string]any `json:"signals,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Provider is implemented by local LLM adapters, SurrealML adapters, and optional cloud adapters.
type Provider interface {
	Name() string
	Supports(task string) bool
	Invoke(ctx context.Context, req Request) (*Response, error)
}

// Router chooses a provider for a task.
type Router struct {
	providers []Provider
}

func New(providers ...Provider) *Router {
	return &Router{providers: providers}
}

func (r *Router) Route(ctx context.Context, req Request) (*Response, error) {
	for _, provider := range r.providers {
		if !req.AllowCloud && provider.Name() == ProviderCloudOpenAICompatible {
			continue
		}

		if provider.Supports(req.Task) {
			return provider.Invoke(ctx, req)
		}
	}

	return nil, errors.New("no model provider supports requested task")
}
