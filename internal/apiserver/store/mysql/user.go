package mysql

import (
	"context"
	v1 "github.com/tiandh987/SharkAgent/api/apiserver/v1"
	metav1 "github.com/tiandh987/SharkAgent/pkg/meta/v1"
	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{ds.db}
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return u.db.Create(&user).Error
}
