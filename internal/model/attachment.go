package model

import (
	"fmt"
	"strings"
)

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Name     string `json:"name"`
	Mime     string `json:"mime"`
	Content  []byte `json:"content"`
	Checksum string `json:"checksum"`
}

func (a Attachment) Validate() error {
	if a.ID == "" || a.RecordID == "" || strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("attachment identity is required")
	}
	if len(a.Content) == 0 {
		return fmt.Errorf("attachment content is required")
	}
	if a.Mime == "" {
		return fmt.Errorf("attachment mime is required")
	}
	return nil
}

func (a Attachment) Size() int { return len(a.Content) }

func (a Attachment) IsImage() bool {
	return strings.HasPrefix(strings.ToLower(a.Mime), "image/")
}
