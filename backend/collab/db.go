package collab

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConn struct {
	pool *pgxpool.Pool

	opTimeout time.Duration
}

func CreateDbConn(opTimeoutTime time.Duration) (*DbConn, error) {
	pool := new(DbConn)

	var err error
	pool.opTimeout = opTimeoutTime
	ctx, cancel := context.WithTimeout(context.Background(), pool.opTimeout * time.Second)
	defer cancel()

	pool.pool, err = pgxpool.New(ctx, os.Getenv("CODELABORATE_DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	
	err = pool.pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	pool.pool.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS Users (
			userID 		char(6),
			name		text,
			email		text,

			CONSTRAINT PRIMARY KEY userID,
			CONSTRAINT UNIQUE NOT NULL name,
			CONSTRAINT UNIQUE NOT NULL email
		)`,
	)

	pool.pool.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS SessionID (
			userID char(6),
			sessionID serial,

			CONSTRAINT PRIMARY KEY userID,
			CONSTRAINT UNIQUE serial
		)`,
	)

	pool.pool.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS Ops (
			userID 		char(6),
			version 	bigserial,
			deleteLen 	integer,
			insertLen	integer,
			insertStr	text,

			CONSTRAINT PRIMARY KEY (userID, version),
			CONSTRAINT FOREIGN KEY (userID) REFERENCES Users(userID) ON DELETE CASCADE,
			CONSTRAINT NOT NULL deleteLen,
			CONSTRAINT NOT NULL insertLen
		)`,
	)



	return pool, nil
}

func (pool *DbConn) CloseConn() error {
	if pool.pool == nil {
		return errors.New("Pool not initialize")
	}
	pool.pool.Close()
	return nil
}

func (pool *DbConn) AddSession(userID string) {

}

func (pool *DbConn) InsertMsg(msg *UpdateMsg) error {
	if pool.pool == nil {
		return errors.New("Pool not initialize")
	}

	ctx, cancel := context.WithTimeout(context.Background(), pool.opTimeout)
	defer cancel()



	return nil
}
