package sshcontrol

import "time"

// Command represents a governed SSH fallback operation.
type Command struct {
	Capability string            `json:"capability"`
	HostID     string            `json:"hostId"`
	Command    string            `json:"command"`
	Args       []string          `json:"args,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	Timeout    time.Duration     `json:"timeout,omitempty"`
	Reason     string            `json:"reason,omitempty"`
}

// Result captures the complete outcome of an SSH fallback command.
type Result struct {
	Capability string    `json:"capability"`
	HostID     string    `json:"hostId"`
	Command    string    `json:"command"`
	Stdout     string    `json:"stdout,omitempty"`
	Stderr     string    `json:"stderr,omitempty"`
	ExitCode   int       `json:"exitCode"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
}
