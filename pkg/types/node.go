package types

import "time"

type Node struct {
	ID       string    `json:"id"`
	IP       string    `json:"ip"`
	CPU      int       `json:"cpu"`
	Memory   int       `json:"memory"`
	LastSeen time.Time `json:"last_seen"`
	Status   string    `json:"status"`
}
