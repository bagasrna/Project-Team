package ikan

import "gorm.io/gorm"

type Repository interface {
	TambahIkan(ikan Ikan) (Ikan, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) TambahIkan(ikan Ikan) (Ikan, error) {
	
	err := r.db.Create(&ikan).Error
	return ikan, err
}