package store

import (
	"gorm.io/gorm"
)

type Factory interface {
	GetDb() *gorm.DB
}
