package models

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	Name   string `json:"name"`
	Lat    string `json:"lat"`
	Lng    string `json:"lng"`
	UserID int64  `json:"user_id"`
}

func (l *Location) Validate() error {
	var existingLocation Location
	// check if user already has these coords
	err := GetDB().Where("user_id = ? AND lat = ? AND lng = ?", l.UserID, l.Lat, l.Lng).First(&existingLocation).Error
	if err == nil {
		return fmt.Errorf("dublicate coords: %#v", l)
	}
	return nil
}

func (l *Location) Create() error {
	if err := l.Validate(); err != nil {
		return err
	}
	if err := GetDB().Create(l).Error; err != nil {
		slog.Error("create location", "error", err)
		return err
	}
	return nil
}

func (l *Location) Update() error {
	err := GetDB().Save(l).Error
	if err != nil {
		slog.Error("gorm save", "error", err)
		return err
	}
	return nil
}

func (l *Location) Delete() error {
	err := GetDB().Unscoped().Delete(&Location{}, l.ID).Error
	if err != nil {
		slog.Error("gorm delete", "error", err)
		return err
	}
	return nil
}
