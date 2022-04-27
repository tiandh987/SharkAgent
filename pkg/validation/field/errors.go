package field

import (
	"fmt"
	utilerrors "github.com/tiandh987/SharkAgent/pkg/errors"

	"github.com/tiandh987/SharkAgent/pkg/util/sets"
)

// TODO: These values are duplicated in api/types.go, but there's a circular dep.  Fix it.
const (
	// ErrorTypeInvalid is used to report malformed values (e.g. failed regex
	// match, too long, out of bounds).  See Invalid().
	ErrorTypeInvalid ErrorType = "FieldValueInvalid"
)

// Error is an implementation of the 'error' interface, which represents a
// field-level validation error.
type Error struct {
	Type     ErrorType
	Field    string
	BadValue interface{}
	Detail   string
}

// ErrorType is a machine readable value providing more detail about why
// a field is invalid.  These values are expected to match 1-1 with
// CauseType in api/types.go.
type ErrorType string

// ErrorList holds a set of Errors.  It is plausible that we might one day have
// non-field errors in this same umbrella package, but for now we don't, so
// we can keep it simple and leave ErrorList here.
type ErrorList []*Error

// ToAggregate converts the ErrorList into an errors.Aggregate.
func (list ErrorList) ToAggregate() utilerrors.Aggregate {
	errs := make([]error, 0, len(list))
	errorMsgs := sets.NewString()
	for _, err := range list {
		msg := fmt.Sprintf("%v", err)
		if errorMsgs.Has(msg) {
			continue
		}
		errorMsgs.Insert(msg)
		errs = append(errs, err)
	}
	return utilerrors.NewAggregate(errs)
}

// Invalid returns a *Error indicating "invalid value".  This is used
// to report malformed values (e.g. failed regex match, too long, out of bounds).
func Invalid(field *Path, value interface{}, detail string) *Error {
	return &Error{ErrorTypeInvalid, field.String(), value, detail}
}


