package database

import (
	"database/sql"
	"log"
	"log/slog"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func New(connectStorage string) (*Storage, error) {
	db, err := sql.Open("postgres", connectStorage)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	_, err = db.Exec(`
		SELECT 'CREATE DATABASE db_avito' 
		WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'db_avito')`)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	database := &Storage{DB: db}

	if err = database.createUsersTable(); err != nil {
		log.Fatal("Ошибка при создании таблицы пользователей: ", err)
	}

	if err = database.createRoomsTable(); err != nil {
		log.Fatal("Ошибка при создании таблицы переговорок: ", err)
	}

	if err = database.createScheduleTable(); err != nil {
		log.Fatal("Ошибка при создании таблицы расписания: ", err)
	}

	if err = database.createSlotsTable(); err != nil {
		log.Fatal("Ошибка при создании таблицы слотов: ", err)
	}

	return database, nil
}

func (storage *Storage) createUsersTable() error {
	_, err := storage.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id UUID PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(10) NOT NULL,
			created_at TIMESTAMP
		)
	`)
	return err
}

func (storage *Storage) createRoomsTable() error {
	_, err := storage.DB.Exec(`
		CREATE TABLE IF NOT EXISTS rooms (
			room_id UUID PRIMARY KEY,
			user_id UUID,
			name VARCHAR(255) UNIQUE NOT NULL,
			description VARCHAR(255),
			capacity INT,
			created_at TIMESTAMP
		)
	`)
	return err
}

func (storage *Storage) createScheduleTable() error {
	_, err := storage.DB.Exec(`
		CREATE TABLE IF NOT EXISTS schedule (
			schedule_id UUID PRIMARY KEY,
			room_id UUID NOT NULL,
			days_of_week INT[] NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL
		)
	`)
	return err
}

func (storage *Storage) createSlotsTable() error {
	_, err := storage.DB.Exec(`
		CREATE TABLE IF NOT EXISTS slots (
			slot_id UUID PRIMARY KEY,
			room_id UUID NOT NULL,
			start_slot TIMESTAMPTZ NOT NULL,
			end_slot TIMESTAMPTZ NOT NULL,
			booking_id UUID,
			booking_status VARCHAR(255),
			booking_user_id UUID,
			booking_created_at TIMESTAMP
		)
	`)
	return err
}
