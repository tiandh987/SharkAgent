package user

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	v1 "github.com/tiandh987/SharkAgent/api/apiserver/v1"
	"github.com/tiandh987/SharkAgent/internal/pkg/code"
	"github.com/tiandh987/SharkAgent/pkg/core"
	"github.com/tiandh987/SharkAgent/pkg/log"
	metav1 "github.com/tiandh987/SharkAgent/pkg/meta/v1"
)

// Create add new user to the storage.
func (u *UserController) Create(c *gin.Context) {
	log.L(c).Info("user create function called.")

	var r v1.User

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	if errs := r.Validate(); len(errs) != 0 {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), nil)

		return
	}

	r.Password, _ = auth.Encrypt(r.Password)
	r.Status = 1

	// Insert the user to the storage.
	if err := u.srv.Users().Create(c, &r, metav1.CreateOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, r)