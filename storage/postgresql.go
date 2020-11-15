package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
	"users-api/models"
)

var ErrNoMatch = fmt.Errorf("no matching record")

type postgreSqlStorage struct {
	pool *pgxpool.Pool
}

type Config struct {
	Username     string
	Password     string
	DatabaseName string
	Host         string
	Port         int
}

const maxPoolConns = 10

func NewPostgreSqlStorage(c Config) (*postgreSqlStorage, error) {
	var pool *pgxpool.Pool
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable pool_max_conns=%d",
		c.Username, c.Password, c.Host, c.Port, c.DatabaseName, maxPoolConns)
	err := backoff.Retry(func() (err error) {
		pool, err = pgxpool.Connect(context.Background(), dsn)
		return
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to database")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to database")
	}

	err = migrateDatabase(conn.Conn())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to migrate database")
	}
	conn.Release()

	return &postgreSqlStorage{
		pool: pool,
	}, nil
}

func migrateDatabase(conn *pgx.Conn) error {
	migrator, err := migrate.NewMigrator(context.Background(), conn, "schema_version")
	if err != nil {
		return errors.Wrapf(err, "unable to create a migrator")
	}

	err = migrator.LoadMigrations("./storage/migrations")
	if err != nil {
		return errors.Wrapf(err, "unable to load migrations")
	}

	err = migrator.Migrate(context.Background())
	if err != nil {
		return errors.Wrapf(err, "unable to migrate")
	}

	ver, err := migrator.GetCurrentVersion(context.Background())
	if err != nil {
		return errors.Wrapf(err, "unable to get current schema version")
	}

	logrus.Infof("Migration done. Current schema version: %v", ver)
	return nil
}

func (p postgreSqlStorage) Get(id int) (*models.User, error) {
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to database")
	}
	defer conn.Release()

	user := models.User{}
	row := conn.QueryRow(context.Background(), `SELECT * FROM users WHERE id = $1;`, id)
	switch err := row.Scan(&user.ID, &user.First, &user.Last, &user.CreatedAt); err {
	case sql.ErrNoRows:
		return &user, ErrNoMatch
	default:
		return &user, err
	}
}

func (p postgreSqlStorage) GetAll() ([]*models.User, error) {
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to database")
	}
	defer conn.Release()

	users := make([]*models.User, 0)
	rows, err := conn.Query(context.Background(), "SELECT * FROM users ORDER BY ID ASC")
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.First, &user.Last, &user.CreatedAt)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (p postgreSqlStorage) Store(user *models.User) (id int, err error) {
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		return -1, errors.Wrapf(err, "unable to connect to database")
	}
	defer conn.Release()

	var createdAt time.Time

	err = conn.QueryRow(
		context.Background(),
		`INSERT INTO users (first, last) VALUES ($1, $2) RETURNING id, created_at`,
		user.First,
		user.Last,
	).Scan(&id, &createdAt)

	if err != nil {
		return
	}

	user.ID = id
	user.CreatedAt = createdAt

	return
}

func (p postgreSqlStorage) Update(id int, userData *models.User) (*models.User, error) {
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to database")
	}
	defer conn.Release()

	user := models.User{}

	err = conn.QueryRow(
		context.Background(),
		`UPDATE users SET first=$1, last=$2 WHERE id=$3 RETURNING id, first, last, created_at;`,
		userData.First,
		userData.Last,
		id,
	).Scan(&user.ID, &user.First, &user.Last, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return &user, ErrNoMatch
		}
		return &user, err
	}

	return &user, nil
}

func (p postgreSqlStorage) Delete(id int) error {
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		return errors.Wrapf(err, "unable to connect to database")
	}
	defer conn.Release()

	_, err = conn.Exec(
		context.Background(),
		`DELETE FROM users WHERE id = $1;`,
		id,
	)

	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (p postgreSqlStorage) UserExist(user *models.User) (bool, error) {
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		return false, errors.Wrapf(err, "unable to connect to database")
	}
	defer conn.Release()

	var exists bool
	row := conn.QueryRow(
		context.Background(),
		`SELECT exists(SELECT * FROM users WHERE first = $1 AND last = $2);`,
		user.First,
		user.Last,
	)

	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
