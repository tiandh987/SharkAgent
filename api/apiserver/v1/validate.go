package v1

import (
	"github.com/tiandh987/SharkAgent/pkg/validation"
	"github.com/tiandh987/SharkAgent/pkg/validation/field"
)

// Validate validates that a user object is valid.
func (u *User) Validate() field.ErrorList {
	val := validation.NewValidator(u)
	allErrs := val.Validate()

	if err := validation.IsValidPassword(u.Password); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("password"), err.Error(), ""))
	}

	return allErrs
}