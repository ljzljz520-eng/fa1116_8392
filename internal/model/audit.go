package model

import "fmt"

type AuditEvent struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Action    string `json:"action"`
	Actor     string `json:"actor"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
	From      string `json:"from"`
	To        string `json:"to"`
}

func (e AuditEvent) Validate() error {
	if e.ID == "" || e.RecordID == "" || e.Action == "" || e.Actor == "" {
		return fmt.Errorf("audit identity and action are required")
	}
	if e.To == "" && e.From == "" {
		return fmt.Errorf("audit transition is required")
	}
	return nil
}

func (e AuditEvent) IsTransition() bool {
	return e.From != "" && e.To != "" && e.From != e.To
}

func (e AuditEvent) Label() string {
	if e.IsTransition() {
		return fmt.Sprintf("%s %s -> %s", e.Action, e.From, e.To)
	}
	return e.Action
}
