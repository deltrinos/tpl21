package db

import (
	"github.com/deltrinos/tpl21/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DB struct {
	Gorm *gorm.DB
}

func Default(connType, connStr string) *DB {
	log.Debug().Msgf("try to connect to database %s...", connType)
	db, err := gorm.Open(connType, connStr)
	if err != nil {
		log.Error().Err(err).Msgf("can't connect to database %v", err)
	} else {
		log.Info().Msgf("connected to database %v", connType)
	}

	return &DB{
		Gorm: db,
	}
}
