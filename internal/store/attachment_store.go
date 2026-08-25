package store

import "gestureparticles/internal/model"

const attachmentBucket = "attachments"

func (d *Database) SaveAttachment(attachment model.Attachment) error {
	if err := attachment.Validate(); err != nil {
		return err
	}
	return d.Put(attachmentBucket, attachment.ID, attachment)
}

func (d *Database) ListAttachments(recordID string) ([]model.Attachment, error) {
	values, err := d.List(attachmentBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.Attachment, 0, len(values))
	for _, data := range values {
		var attachment model.Attachment
		if err := unmarshal(data, &attachment); err != nil {
			return nil, err
		}
		if recordID == "" || attachment.RecordID == recordID {
			result = append(result, attachment)
		}
	}
	return result, nil
}
