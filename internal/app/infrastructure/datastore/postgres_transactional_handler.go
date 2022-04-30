package datastore

import (
	"context"
	"errors"
	"fmt"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/repository/basedbhandler"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"time"
)

type PostgresHandlerTX struct {
	pool *pgxpool.Pool
	log  *infrastructure.Logger
}

func NewPostgresHandlerTX(ctx context.Context, dataSource string, log *infrastructure.Logger) (*PostgresHandlerTX, error) {
	poolConfig, err := pgxpool.ParseConfig(dataSource)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = 5
	poolConfig.MinConns = 2
	poolConfig.MaxConnIdleTime = time.Second * 120
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	postgresHandler := new(PostgresHandlerTX)
	postgresHandler.pool = pool
	postgresHandler.log = log
	return postgresHandler, nil
}

func (handler *PostgresHandlerTX) NewTx(ctx context.Context) (pgx.Tx, error) {
	return handler.pool.Begin(ctx)
}

func (handler *PostgresHandlerTX) getTx(ctx context.Context) (tx pgx.Tx, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("can't get tx: conversion error")
			handler.log.Error("PostgresHandlerTX: can't get tx", zap.Error(err))
		}
	}()
	ctxValue := ctx.Value(basedbhandler.TransactionKey("tx"))
	if ctxValue == nil {
		return nil, errors.New("can't get tx: nil value got")
	}
	tx = ctxValue.(pgx.Tx)
	return tx, err
}

func (handler *PostgresHandlerTX) Commit(ctx context.Context) error {
	tx, err := handler.getTx(ctx)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		handler.log.Error("Can't commit transaction", zap.Error(err))
		return err
	}
	return err
}

func (handler *PostgresHandlerTX) Rollback(ctx context.Context) error {
	tx, err := handler.getTx(ctx)
	if err != nil {
		return err
	}
	err = tx.Rollback(ctx)
	if err != nil {
		handler.log.Error("Can't commit transaction", zap.Error(err))
		return err
	}
	return err
}

func (handler *PostgresHandlerTX) Execute(ctx context.Context, statement string, args ...interface{}) error {
	tx, err := handler.getTx(ctx)
	if err == nil {
		if len(args) > 0 {
			_, err = tx.Exec(ctx, statement, args...)
		} else {
			_, err = tx.Exec(ctx, statement)
		}
	} else {
		conn, e := handler.pool.Acquire(ctx)
		if e != nil {
			return e
		}
		defer conn.Release()

		if len(args) > 0 {
			_, e = conn.Exec(ctx, statement, args...)
		} else {
			_, e = conn.Exec(ctx, statement)
		}
		err = e
	}
	return err
}

func (handler *PostgresHandlerTX) ExecuteBatch(ctx context.Context, statement string, args [][]interface{}) error {
	var (
		err error
		ct  pgconn.CommandTag
		br  pgx.BatchResults
	)

	batch := &pgx.Batch{}
	if len(args) > 0 {
		for _, argset := range args {
			batch.Queue(statement, argset...)
		}
	} else {
		return nil
	}
	tx, err := handler.getTx(ctx)

	if err == nil {
		br = tx.SendBatch(context.Background(), batch)
	} else {
		conn, err := handler.pool.Acquire(ctx)
		if err != nil {
			return err
		}
		defer conn.Release()
		br = conn.SendBatch(context.Background(), batch)
	}
	ct, err = br.Exec()

	if err != nil {
		return err
	}
	fmt.Println(ct.RowsAffected())
	return nil
}

func (handler *PostgresHandlerTX) QueryRow(ctx context.Context, statement string, args ...interface{}) (basedbhandler.Row, error) {
	var row pgx.Row
	tx, err := handler.getTx(ctx)
	if err == nil {
		if len(args) > 0 {
			row = tx.QueryRow(ctx, statement, args...)
		} else {
			row = tx.QueryRow(ctx, statement)
		}
	} else {
		conn, err := handler.pool.Acquire(ctx)
		if err != nil {
			return nil, err
		}
		defer conn.Release()
		if len(args) > 0 {
			row = conn.QueryRow(ctx, statement, args...)
		} else {
			row = conn.QueryRow(ctx, statement)
		}
	}
	return row, nil
}

func (handler *PostgresHandlerTX) Query(ctx context.Context, statement string, args ...interface{}) (basedbhandler.Rows, error) {
	var rows pgx.Rows
	tx, err := handler.getTx(ctx)
	if err == nil {
		if len(args) > 0 {
			rows, err = tx.Query(ctx, statement, args...)
		} else {
			rows, err = tx.Query(ctx, statement)
		}
	} else {
		conn, e := handler.pool.Acquire(ctx)
		if e != nil {
			return nil, e
		}
		defer conn.Release()
		if len(args) > 0 {
			rows, e = conn.Query(ctx, statement, args...)
		} else {
			rows, e = conn.Query(ctx, statement)
		}
		err = e
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (handler *PostgresHandlerTX) Close() {
	if handler != nil {
		handler.pool.Close()
	}
}
