package store

import (
	"fmt"
	"strings"
)

var (
	ErrRecordNotFound = &RecordNotFound{}
	ErrFieldMissing   = &RequiredFieldMissing{}
)

type RequiredFieldMissing struct {
	fields []string
}

func (e *RequiredFieldMissing) AddFields(fields ...string) *RequiredFieldMissing {
	e.fields = append(e.fields, fields...)
	return e
}

func (e *RequiredFieldMissing) Error() string {
	return fmt.Sprintf("fields missing: %s", strings.Join(e.fields, ", "))
}

type RecordNotFound struct {
	record string
}

func (e *RecordNotFound) Record(record string) *RecordNotFound {
	e.record = record
	return e
}

func (e *RecordNotFound) Error() string {
	return fmt.Sprintf("record %s not found", e.record)
}
