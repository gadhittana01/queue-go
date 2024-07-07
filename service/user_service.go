package service

import (
	"context"

	querier "github.com/gadhittana-01/queue-go/db/repository"
	"github.com/gadhittana-01/queue-go/dto"
	"github.com/gadhittana-01/queue-go/utils"
	"github.com/jackc/pgx/v5"
)

const (
	FailedToParseDate     = "Failed to parse date"
	FailedToCreateUser    = "Failed to create user"
	EmailAlreadyExist     = "Email already exist"
	FailedToFindUser      = "Failed to find user"
	WrongCredentials      = "Wrong credentials"
	FailedToGenerateToken = "Failed to generate token"
	UserNotFound          = "User not found"
)

type UserSvc interface {
	SignUp(ctx context.Context, input dto.SignUpReq) dto.SignUpRes
	SignIn(ctx context.Context, input dto.SignInReq) dto.SignInRes
}

type UserSvcImpl struct {
	repo   querier.Repository
	config *utils.BaseConfig
	token  utils.TokenClient
}

func NewUserSvc(
	repo querier.Repository,
	config *utils.BaseConfig,
	token utils.TokenClient,
) UserSvc {
	return &UserSvcImpl{
		repo:   repo,
		config: config,
		token:  token,
	}
}

func (s *UserSvcImpl) SignUp(ctx context.Context, input dto.SignUpReq) dto.SignUpRes {
	var resp dto.SignUpRes
	var user querier.User
	var token utils.GenerateTokenResp
	var err error

	err = utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		isExists, err := repoTx.CheckEmailExists(ctx, input.Email)
		if isExists {
			return utils.CustomErrorWithTrace(err, EmailAlreadyExist, 422)
		}

		user, err = repoTx.CreateUser(ctx, querier.CreateUserParams{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
		})
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToCreateUser, 422)
		}

		token, err = s.token.GenerateToken(utils.GenerateTokenReq{
			UserID: user.ID.String(),
		})
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToGenerateToken, 400)
		}

		return nil
	})
	utils.PanicIfError(err)

	resp = dto.SignUpRes{
		ID:       user.ID.String(),
		Name:     user.Name,
		Token:    token.Token,
		ExpToken: token.ExpToken,
	}

	return resp
}

func (s *UserSvcImpl) SignIn(ctx context.Context, input dto.SignInReq) dto.SignInRes {
	var resp dto.SignInRes
	var user querier.User
	var err error

	user, err = s.repo.FindUserByEmail(ctx, input.Email)
	utils.PanicIfAppError(err, FailedToFindUser, 400)

	if user.Password != input.Password {
		utils.PanicIfAppError(err, WrongCredentials, 400)
	}

	token, err := s.token.GenerateToken(utils.GenerateTokenReq{
		UserID: user.ID.String(),
	})
	if err != nil {
		utils.PanicIfAppError(err, FailedToGenerateToken, 400)
	}

	resp = dto.SignInRes{
		ID:       user.ID.String(),
		Name:     user.Name,
		Token:    token.Token,
		ExpToken: token.ExpToken,
	}

	return resp
}
