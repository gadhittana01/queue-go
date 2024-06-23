package service

import (
	"context"
	"time"

	"github.com/gadhittana-01/form-go/constant"
	querier "github.com/gadhittana-01/form-go/db/repository"
	"github.com/gadhittana-01/form-go/dto"
	"github.com/gadhittana-01/form-go/utils"
	"github.com/jackc/pgx/v5"
)

const (
	FailedToParseDate  = "Failed to parse date"
	FailedToCreateUser = "Failed to create user"
)

type UserSvc interface {
	CreateUser(ctx context.Context, input dto.CreateUserReq) dto.CreateUserRes
}

type UserSvcImpl struct {
	repo   querier.Repository
	config *utils.BaseConfig
}

func NewUserSvc(
	repo querier.Repository,
	config *utils.BaseConfig,
) UserSvc {
	return &UserSvcImpl{
		repo:   repo,
		config: config,
	}
}

func (s *UserSvcImpl) CreateUser(ctx context.Context, input dto.CreateUserReq) dto.CreateUserRes {
	var resp dto.CreateUserRes
	var user querier.User

	err := utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		date, err := time.Parse(constant.DateFormat, input.DateOfBirth)
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToParseDate, 400)
		}

		user, err = repoTx.CreateUser(ctx, querier.CreateUserParams{
			Name:           input.Name,
			IdentityNumber: input.IdentityNumber,
			Email:          input.Email,
			DateOfBirth:    date,
		})
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToCreateUser, 422)
		}

		return nil
	})
	utils.PanicIfError(err)

	resp = dto.CreateUserRes{
		ID:             user.ID.String(),
		Name:           user.Name,
		IdentityNumber: user.IdentityNumber,
		Email:          user.Email,
		DateOfBirth:    user.DateOfBirth.Format(constant.DateFormat),
	}

	return resp
}
