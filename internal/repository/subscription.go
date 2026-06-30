package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"subscription-service/internal/model"
)

type Repository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(subscription model.Subscription) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO subscriptions
		(id, service_name, price, user_id, start_date, end_date)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	)

	return err
}

func (r *Repository) GetAll() ([]model.Subscription, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, service_name, price, user_id, start_date, end_date
		 FROM subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription

	for rows.Next() {
		var s model.Subscription

		err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
		)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, s)
	}

	return subscriptions, nil
}

func (r *Repository) GetByID(id string) (*model.Subscription, error) {
	var subscription model.Subscription

	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1`,
		id,
	).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
	)

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(
		context.Background(),
		`DELETE FROM subscriptions WHERE id = $1`,
		id,
	)

	return err
}

func (r *Repository) Update(subscription model.Subscription) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE subscriptions
		SET service_name = $1,
		    price = $2,
		    user_id = $3,
		    start_date = $4,
		    end_date = $5
		WHERE id = $6`,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.ID,
	)

	return err
}

func (r *Repository) CalculateCost(userID, serviceName, from, to string) (int, error) {
	var total int

	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE 1=1`
	args := []interface{}{}
	arg := 1

	if userID != "" {
		query += fmt.Sprintf(" AND user_id = $%d::uuid", arg)
		args = append(args, userID)
		arg++
	}

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", arg)
		args = append(args, serviceName)
		arg++
	}

	if from != "" {
		query += fmt.Sprintf(" AND start_date >= $%d::date", arg)
		args = append(args, from+"-01")
		arg++
	}

	if to != "" {
		query += fmt.Sprintf(" AND start_date <= ($%d::date + interval '1 month' - interval '1 day')", arg)
		args = append(args, to+"-01")
	}

	err := r.db.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}
