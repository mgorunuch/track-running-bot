package main

import "database/sql"

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

	rows, err := db.Query(`select telegram_id, sum(distance) as distance, (select name from data dt where dt.telegram_id = data.telegram_id order by created_at desc limit 1) as name from data group by telegram_id`)
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

func insertData(db *sql.DB, id int, distance float64, name string) error {
	_, err := db.Exec(`insert into data (telegram_id, distance, name) values ($1, $2, $3);`, id, distance, name)
	if err != nil {
		return err
	}

	return nil
}
