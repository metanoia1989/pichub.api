package repository

import (
	"pichub.api/infra/database"
	"pichub.api/infra/logger"
)

func Save(model interface{}) error {
	if err := database.DB.Create(model).Error; err != nil {
		logger.Errorf("error, not save data %v", err)
		return err
	}
	return nil
}

func Get(model interface{}, query interface{}, args ...interface{}) error {
	if err := database.DB.Where(query, args...).Find(model).Error; err != nil {
		logger.Errorf("error getting data: %v", err)
		return err
	}
	return nil
}

func GetOne(model interface{}, query interface{}, args ...interface{}) error {
	if err := database.DB.Where(query, args...).Last(model).Error; err != nil {
		logger.Errorf("error getting single record: %v", err)
		return err
	}
	return nil
}

func Update(model interface{}, query interface{}, args ...interface{}) error {
	if err := database.DB.Where(query, args...).Updates(model).Error; err != nil {
		logger.Errorf("error updating data: %v", err)
		return err
	}
	return nil
}
