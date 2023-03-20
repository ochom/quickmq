package dto

import "time"

// PublishRequest is a struct that holds the data for a publish request
type PublishRequest struct {
	Queue string        `json:"queue" binding:"required"`
	Data  []byte        `json:"data" binding:"required"`
	Delay time.Duration `json:"delay" binding:"required"`
}
