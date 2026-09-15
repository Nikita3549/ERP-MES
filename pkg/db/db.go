// Package db provides database connectivity
package db

import (
	"context"

	"erp-mes/internal/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(conf *configs.Config) (*DB, error) {
	db, err := gorm.Open(postgres.Open(conf.DBConfig.DSN), &gorm.Config{})
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
