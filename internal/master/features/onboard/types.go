package onboard

import "time"

type CreateTokenResponse struct {
	Token     string
	ExpiresAt time.Time
}

type RegisterNodeRequest struct {
	Token    string
	NodeID   string
	Hostname string
	Name     string
	Region   string
	Type     string
	CPU      int64
	Memory   int64
	Arch     string
	IP       string
}

type RegisterNodeResponse struct {
	AuthKey  string
	NodeName string
	Labels   []string
}
