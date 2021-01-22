package config

import (
	"context"
)

// Endpoint to listen for requests and what to execute
type Endpoint struct {
	Command                string   `json:"command"`
	Arguments              []string `json:"arguments"`
	HTTPMethod             string   `json:"http_method"`
	AllowExternalArguments bool     `json:"allow_external_arguments"`
	AllowStdin             bool     `json:"allow_stdin"`
	DiscardStderr          bool     `json:"discard_stderr"`
	DiscardStdout          bool     `json:"discard_stdout"`
}

// Config information for httpwrap
type Config struct {
	Endpoints     map[string]*Endpoint `json:"endpoints"`
	ListenAddress string               `json:"listen_address"`
	ListenPort    uint16               `json:"listen_port"`
}

// Source that can get a Config instance
type Source interface {
	Config(context.Context) (*Config, error)
}
