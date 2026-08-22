package ssehub

import "errors"

var (
	ErrClosed   = errors.New("ssehub: closed")
	ErrInvalid  = errors.New("ssehub: invalid")
	ErrNotFound = errors.New("ssehub: not found")
	ErrConflict = errors.New("ssehub: conflict")
	// ErrAudit is returned when an acknowledged write could not be made
	// durable on the audit side (closed handle, I/O error, etc.). Callers
	// that ask for an acked write must treat this as a failure and not
	// advance their pipeline as if the record had been audited.
	ErrAudit = errors.New("ssehub: audit write failed")
)
