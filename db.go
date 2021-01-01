package main

import (
	"database/sql"
	"time"
)

func connect(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func refreshData(db *sql.DB) error {
	store = map[int]float64{}
	names = map[int]string{}

	rows, err := db.Query(`
		select telegram_id,
		       sum(distance) as distance,
		       (select name from data dt where dt.telegram_id = data.telegram_id order by created_at desc limit 1) as name
		from data
		where created_at > $1
		group by telegram_id
	`, startDate)
	if err != nil {
		return err
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	for rows.Next() {
		var (
			tgID int
			dist float64
			name string
		)

		err = rows.Scan(&tgID, &dist, &name)
		if err != nil {
			return err
		}

		store[tgID] = dist
		names[tgID] = name
	}

	return nil
}

type RunningItem struct {
	createdAt time.Time
	distance  float64
}

func getSingleUserData(db *sql.DB, tgID int) ([]RunningItem, error) {
	rows, err := db.Query(`
		select created_at,
		       distance
		from data
		where telegram_id = $1
		  and created_at > $2
	`, tgID, startDate)
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	dt := make([]RunningItem, 0)

	for rows.Next() {
		var (
			dist      float64
			createdAt time.Time
		)

		err = rows.Scan(&createdAt, &dist)
		if err != nil {
			return nil, err
		}

		dt = append(dt, RunningItem{
			createdAt: createdAt,
			distance:  dist,
		})
	}

	return dt, nil
}

func insertData(db *sql.DB, id int, distance float64, name string) error {
	_, err := db.Exec(`insert into data (telegram_id, distance, name) values ($1, $2, $3);`, id, distance, name)
	if err != nil {
		return err
	}

	return nil
}
