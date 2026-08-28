package store

import (
	"fmt"
	"gestureparticles/internal/model"
	"sort"
	"strings"
)

const recordBucket = "records"

func (d *Database) SaveRecord(record model.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return d.Put(recordBucket, record.ID, record)
}

func (d *Database) LoadRecord(id string) (model.Record, error) {
	var record model.Record
	found, err := d.Get(recordBucket, id, &record)
	if err != nil {
		return record, err
	}
	if !found {
		return record, fmt.Errorf("record %s not found", id)
	}
	return record, nil
}

func (d *Database) DeleteRecord(id string) error { return d.Delete(recordBucket, id) }

func (d *Database) ListRecords(classroom string) ([]model.Record, error) {
	values, err := d.List(recordBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.Record, 0, len(values))
	for _, data := range values {
		var record model.Record
		if err := unmarshal(data, &record); err != nil {
			return nil, err
		}
		if classroom == "" || strings.EqualFold(record.Classroom, classroom) {
			result = append(result, record)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Sequence == result[j].Sequence {
			return result[i].ID < result[j].ID
		}
		return result[i].Sequence < result[j].Sequence
	})
	return result, nil
}

func (d *Database) NextSequence(classroom string) (int, error) {
	records, err := d.ListRecords(classroom)
	if err != nil {
		return 0, err
	}
	next := 1
	for _, record := range records {
		if record.Sequence >= next {
			next = record.Sequence + 1
		}
	}
	return next, nil
}
