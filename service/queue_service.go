package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gadhittana-01/queue-go/constant"
	querier "github.com/gadhittana-01/queue-go/db/repository"
	"github.com/gadhittana-01/queue-go/dto"
	"github.com/gadhittana-01/queue-go/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
)

const (
	FailedToFindQueue         = "Failed to find queue"
	FailedToParseStringToUUID = "Failed to parse string to uuid"
	FailedToFindLatestQueue   = "Failed to find latest queue"
	FailedToCreateQueue       = "Failed to create queue"
	QueueNotFound             = "Queue not found"
	QueueHasBeenProcess       = "Queue has been process"
)

type QueueSvc interface {
	GetQueue(ctx context.Context) []dto.GetQueueRes
	CreateQueue(ctx context.Context) dto.CreateQueueRes
	ProcessQueue(ctx context.Context, input dto.ProcessQueueReq) dto.ProcessQueueRes
	DeleteQueue(ctx context.Context, input dto.DeleteQueueReq) dto.SimpleMessageRes
}

type QueueSvcImpl struct {
	repo   querier.Repository
	config *utils.BaseConfig
}

func NewQueueSvc(
	repo querier.Repository,
	config *utils.BaseConfig,
) QueueSvc {
	return &QueueSvcImpl{
		repo:   repo,
		config: config,
	}
}

func formatDuration(duration int64) string {
	d := time.Duration(duration) * time.Microsecond

	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second

	return fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
}

func formatQueueNumber(recentQueueNumber string) string {
	if recentQueueNumber == "" {
		return "A001"
	}

	prefix := recentQueueNumber[:1] // Assuming prefix is always one character

	numStr := recentQueueNumber[1:]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return recentQueueNumber // Return original if unable to parse
	}

	num++

	format := fmt.Sprintf("%s%%0%dd", prefix, len(numStr))
	return fmt.Sprintf(format, num)
}

func (s *QueueSvcImpl) GetQueue(ctx context.Context) []dto.GetQueueRes {
	var err error

	queue, err := s.repo.FindQueue(ctx)
	utils.PanicIfAppError(err, FailedToFindQueue, 400)

	return lo.Map(queue, func(item querier.Queue, index int) dto.GetQueueRes {
		return dto.GetQueueRes{
			ID:               item.ID.String(),
			QueueNumber:      item.QueueNumber,
			UserID:           item.UserID.String(),
			ArrivalTime:      item.ArrivalTime.Format(constant.TimeFormat),
			ServiceTime:      item.ServiceTime.Time.Format(constant.TimeFormat),
			TotalWaitingTime: formatDuration(item.TotalWaitingTime.Microseconds),
			CreatedAt:        item.CreatedAt.Format(constant.TimeFormat),
			UpdatedAt:        item.UpdatedAt.Format(constant.TimeFormat),
		}
	})
}

func (s *QueueSvcImpl) CreateQueue(ctx context.Context) dto.CreateQueueRes {
	var err error
	var result dto.CreateQueueRes
	authPayload := utils.GetRequestCtx(ctx, constant.UserSession)
	now := time.Now()

	err = utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		userID, err := uuid.Parse(authPayload.UserID)
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToParseStringToUUID, 400)
		}

		isExists, err := repoTx.CheckUserExists(ctx, userID)
		if !isExists {
			return utils.CustomError(UserNotFound, 400)
		}

		queue, err := repoTx.FindLatestQueue(ctx)
		if err != nil && err != pgx.ErrNoRows {
			return utils.CustomErrorWithTrace(err, FailedToFindLatestQueue, 400)
		}

		res, err := repoTx.CreateQueue(ctx, querier.CreateQueueParams{
			QueueNumber: formatQueueNumber(queue.QueueNumber),
			UserID:      userID,
			ArrivalTime: now,
		})
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToCreateQueue, 422)
		}

		result = dto.CreateQueueRes{
			ID:          res.ID.String(),
			QueueNumber: res.QueueNumber,
			UserID:      res.UserID.String(),
			ArrivalTime: res.ArrivalTime.Format(constant.TimeFormat),
			CreatedAt:   res.CreatedAt.Format(constant.TimeFormat),
			UpdatedAt:   res.UpdatedAt.Format(constant.TimeFormat),
		}

		return nil
	})
	utils.PanicIfError(err)

	return result
}

func (s *QueueSvcImpl) ProcessQueue(ctx context.Context, input dto.ProcessQueueReq) dto.ProcessQueueRes {
	var err error
	var result dto.ProcessQueueRes
	authPayload := utils.GetRequestCtx(ctx, constant.UserSession)
	now := time.Now()

	err = utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		userID, err := uuid.Parse(authPayload.UserID)
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToParseStringToUUID, 400)
		}

		isExists, err := repoTx.CheckUserExists(ctx, userID)
		if !isExists {
			return utils.CustomError(UserNotFound, 400)
		}

		queue, err := repoTx.FindQueueByID(ctx, querier.FindQueueByIDParams{
			ID:     input.QueueID,
			UserID: userID,
		})
		if err != nil && err != pgx.ErrNoRows {
			return utils.CustomErrorWithTrace(err, FailedToFindLatestQueue, 400)
		}

		if queue == (querier.Queue{}) {
			return utils.CustomError(QueueNotFound, 400)
		}

		if queue.TotalWaitingTime.Valid {
			return utils.CustomError(QueueHasBeenProcess, 400)
		}

		duration := now.Sub(queue.ArrivalTime).Microseconds()

		res, err := repoTx.UpdateQueue(ctx, querier.UpdateQueueParams{
			ID: queue.ID,
			ServiceTime: sql.NullTime{
				Time:  now,
				Valid: true,
			},
			TotalWaitingTime: pgtype.Interval{
				Microseconds: duration,
				Valid:        true,
			},
		})
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToCreateQueue, 422)
		}

		result = dto.ProcessQueueRes{
			QueueNumber:      res.QueueNumber,
			UserID:           res.UserID.String(),
			ServiceTime:      res.ServiceTime.Time.Format(constant.TimeFormat),
			TotalWaitingTime: formatDuration(res.TotalWaitingTime.Microseconds),
			ArrivalTime:      res.ArrivalTime.Format(constant.TimeFormat),
			CreatedAt:        res.CreatedAt.Format(constant.TimeFormat),
			UpdatedAt:        res.UpdatedAt.Format(constant.TimeFormat),
		}

		return nil
	})
	utils.PanicIfError(err)

	return result
}

func (s *QueueSvcImpl) DeleteQueue(ctx context.Context, input dto.DeleteQueueReq) dto.SimpleMessageRes {
	var err error
	var result dto.SimpleMessageRes
	authPayload := utils.GetRequestCtx(ctx, constant.UserSession)

	err = utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		userID, err := uuid.Parse(authPayload.UserID)
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToParseStringToUUID, 400)
		}

		isExists, err := repoTx.CheckUserExists(ctx, userID)
		if !isExists {
			return utils.CustomError(UserNotFound, 400)
		}

		queue, err := repoTx.FindQueueByID(ctx, querier.FindQueueByIDParams{
			ID:     input.QueueID,
			UserID: userID,
		})
		if err != nil && err != pgx.ErrNoRows {
			return utils.CustomErrorWithTrace(err, FailedToFindLatestQueue, 400)
		}

		if queue == (querier.Queue{}) {
			return utils.CustomError(QueueNotFound, 400)
		}

		if queue.TotalWaitingTime.Valid {
			return utils.CustomError(QueueHasBeenProcess, 400)
		}

		err = repoTx.DeleteQueue(ctx, querier.DeleteQueueParams{
			ID:     queue.ID,
			UserID: userID,
		})
		if err != nil {
			return utils.CustomErrorWithTrace(err, FailedToCreateQueue, 422)
		}

		result = dto.SimpleMessageRes{
			Message: fmt.Sprintf("Queue with id %s has been deleted", queue.ID),
		}

		return nil
	})
	utils.PanicIfError(err)

	return result
}
