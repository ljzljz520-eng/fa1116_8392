package store

import (
	"gestureparticles/internal/model"
)

const auditBucket = "audits"

func (d *Database) SaveAudit(event model.AuditEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	return d.Put(auditBucket, event.ID, event)
}

func (d *Database) ListAudits(recordID string) ([]model.AuditEvent, error) {
	values, err := d.List(auditBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.AuditEvent, 0, len(values))
	for _, data := range values {
		var event model.AuditEvent
		if err := unmarshal(data, &event); err != nil {
			return nil, err
		}
		if recordID == "" || event.RecordID == recordID {
			result = append(result, event)
		}
	}
	return result, nil
}
