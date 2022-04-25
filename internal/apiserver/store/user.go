package store

import (
	"context"
	v1 "github.com/tiandh987/SharkAgent/api/apiserver/v1"
	metav1 "github.com/tiandh987/SharkAgent/pkg/meta/v1"
)

// UserStore defines the user storage interface.
type UserStore interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
}
