package repository

import (
	"mini-crm/internal/models"
)

type TokenRepository interface {
	Save(token *models.RefreshToken) error
	Delete(token string) error
	Find(token string) (*models.RefreshToken, error)
}
