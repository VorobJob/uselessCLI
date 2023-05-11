package repository

import (
	"CLITool/internal/models"
)

type Repository interface {
	CreateTable() error
	CreateRecord(user models.User) error
	GetUniqueRecords() error
	AutoFill() error
	GetRecords() error
	CreateIndex() error
	GetIndexedRecords() error
}
