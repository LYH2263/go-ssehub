package ssehub
import "errors"
var (
        ErrClosed = errors.New("ssehub: closed")
        ErrInvalid = errors.New("ssehub: invalid")
        ErrNotFound = errors.New("ssehub: not found")
        ErrConflict = errors.New("ssehub: conflict")
)
