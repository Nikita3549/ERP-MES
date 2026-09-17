// Package db provides database connectivity
package db

import (
	"context"
	"fmt"

	"erp-mes/internal/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(conf *configs.Config) (*DB, error) {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", conf.DBConfig.Host, conf.DBConfig.User, conf.DBConfig.Password, conf.DBConfig.Name, conf.DBConfig.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}

func (d *DB) Health(ctx context.Context) error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
