package v1

import (
	"context"
	"github.com/marmotedu/errors"
	v1 "github.com/tiandh987/SharkAgent/api/apiserver/v1"
	"github.com/tiandh987/SharkAgent/internal/apiserver/store"
	"github.com/tiandh987/SharkAgent/internal/pkg/code"
	metav1 "github.com/tiandh987/SharkAgent/pkg/meta/v1"
	"regexp"
)

type UserSrv interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
}

type userService struct {
	store store.Factory
}

var _ UserSrv = (*userService)(nil)

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

func (u *userService) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	if err := u.store.Users().Create(ctx, user, opts); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'idx_name'", err.Error()); match {
			return errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}

		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

