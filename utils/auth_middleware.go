package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gadhittana-01/queue-go/constant"
)

const (
	HeaderContentType           = "Content-Type"
	HeaderUserAgent             = "User-Agent"
	HeaderAuthorization         = "Authorization"
	HeaderAuthorizationCustomer = "Authorization-Customer"
	HeaderXFAuthorization       = "X-Forwarded-Authorization"
	HeaderXUserID               = "X-User-Id"
)

type AuthPayload struct {
	UserID string `json:"userId"`
}

type AuthMiddleware interface {
	CheckIsAuthenticated(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc
}

type AuthMiddlewareImpl struct {
	config *BaseConfig
	token  TokenClient
}

func NewAuthMiddleware(config *BaseConfig, token TokenClient) AuthMiddleware {
	return &AuthMiddlewareImpl{
		config: config,
		token:  token,
	}
}

func (m *AuthMiddlewareImpl) CheckIsAuthenticated(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(HeaderXFAuthorization)
		if header == "" {
			header = r.Header.Get(HeaderAuthorization)
		}

		if header == "" || !strings.Contains(header, "Bearer ") {
			PanicIfError(CustomError("unauthorized", 401))
		}
		authToken := strings.Split(header, " ")[1]

		res, err := m.token.DecodeToken(DecodeTokenReq{
			Token: authToken,
		})
		if err != nil {
			LogInfo(fmt.Sprintf("failed when decode token, error: %v", err))
			PanicIfError(CustomError("unauthorized", 401))
		}

		authCtx := AppendRequestCtx(r, constant.UserSession, &AuthPayload{
			UserID: res.UserID,
		})
		handler(w, r.WithContext(authCtx))
	}
}

func AppendRequestCtx(r *http.Request, ctxKey string, input interface{}) context.Context {
	return context.WithValue(r.Context(), ctxKey, input)
}

func GetRequestCtx(ctx context.Context, ctxKey string) *AuthPayload {
	ctxVal := ctx.Value(ctxKey)
	if ctxVal != nil {
		return ctxVal.(*AuthPayload)
	}

	return &AuthPayload{}
}
