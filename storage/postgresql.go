package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
	"users-api/models"
)

var ErrNoMatch = fmt.Errorf("no matching record")

type postgreSqlStorage struct {
	conn *sql.DB
}

type Config struct {
	Username     string
	Password     string
	DatabaseName string
	Host         string
	Port         int
}

func NewPostgreSqlStorage(c Config) *postgreSqlStorage {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.Username, c.Password, c.DatabaseName)
	logrus.Info(dsn)
	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		panic(errors.Wrapf(err, "unable to connect to database"))
	}

	err = conn.Ping()
	if err != nil {
		panic(errors.Wrapf(err, "trying to ping database"))
	}

	return &postgreSqlStorage{
		conn: conn,
	}
}

func (p postgreSqlStorage) Get(id int) (*models.User, error) {
	user := models.User{}
	row := p.conn.QueryRow(`SELECT * FROM users WHERE id = $1;`, id)
	switch err := row.Scan(&user.ID, &user.First, &user.Last, &user.CreatedAt); err {
	case sql.ErrNoRows:
		return &user, ErrNoMatch
	default:
		return &user, err
	}
}

func (p postgreSqlStorage) GetAll() ([]*models.User, error) {
	users := make([]*models.User, 0)
	rows, err := p.conn.Query("SELECT * FROM users ORDER BY ID ASC")
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
	var createdAt time.Time

	err = p.conn.QueryRow(
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
	user := models.User{}

	err := p.conn.QueryRow(
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
	_, err := p.conn.Exec(
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
	var exists bool
	row := p.conn.QueryRow(
		`SELECT exists(SELECT * FROM users WHERE first = $1 AND last = $2);`,
		user.First,
		user.Last,
	)

	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
