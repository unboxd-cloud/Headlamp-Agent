package events

import "time"

const (
	EventTypeHeartbeat    = "heartbeat"
	EventTypeNotification = "notification"
	EventTypeIncident     = "incident"
	EventTypeAudit        = "audit"
	EventTypeHealthChange = "health_change"
)

// EventEnvelope is the common record shape stored in SurrealDB event tables
// and streamed through SurrealDB live/event subscriptions.
type EventEnvelope[T any] struct {
	ID          string    `json:"id,omitempty"`
	Kind        string    `json:"kind"`
	Type        string    `json:"type"`
	Source      string    `json:"source"`
	SubjectKind string    `json:"subjectKind,omitempty"`
	SubjectID   string    `json:"subjectId,omitempty"`
	Severity    string    `json:"severity,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	Payload     T         `json:"payload"`
}

type HeartbeatPayload struct {
	AgentID      string            `json:"agentId,omitempty"`
	WatcherID    string            `json:"watcherId,omitempty"`
	HostID       string            `json:"hostId,omitempty"`
	Status       string            `json:"status"`
	Version      string            `json:"version,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	Health       map[string]string `json:"health,omitempty"`
}

type NotificationPayload struct {
	Title              string   `json:"title"`
	Message            string   `json:"message"`
	RecommendedActions []string `json:"recommendedActions,omitempty"`
}

type IncidentPayload struct {
	Title             string   `json:"title"`
	AffectedArtifacts []string `json:"affectedArtifacts,omitempty"`
	RecoveryState     string   `json:"recoveryState,omitempty"`
	NextAction         string   `json:"nextAction,omitempty"`
}

type AuditPayload struct {
	Action     string         `json:"action"`
	Decision   string         `json:"decision,omitempty"`
	Reason     string         `json:"reason,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
}
