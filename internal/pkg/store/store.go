package store

import (
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
)

type Factory interface {
	GetDb() *gorm.DB
	Getclient() *redis.Client
}
