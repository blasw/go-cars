package storage

import (
	"errors"
	"go-cars/storage/models"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CarFilter struct {
	RegNum          string
	Mark            string
	Model           string
	Year            int
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic string
}

// GetFilteredCars retrieves a list of filtered cars from the storage.
//
// The filter parameter is used to specify the criteria for filtering the cars.
// The limit parameter specifies the maximum number of cars to retrieve.
// The offset parameter specifies the number of cars to skip before starting to retrieve.
// The function returns a slice of models.Cars and an error if any.
func (s *Storage) GetFilteredCars(filter *CarFilter, limit int, offset int) ([]models.Cars, error) {
	var cars []models.Cars

	query := s.DB.Model(&models.Cars{}).Preload("Owner").Joins("join owners on owners.id = cars.owner_id")

	if filter != nil {
		if filter.RegNum != "" {
			query = query.Where("lower(reg_num) like ?", "%"+strings.ToLower(filter.RegNum)+"%")
		}
		if filter.Mark != "" {
			query = query.Where("lower(mark) like ?", "%"+strings.ToLower(filter.Mark)+"%")
		}
		if filter.Model != "" {
			query = query.Where("lower(model) like ?", "%"+strings.ToLower(filter.Model)+"%")
		}
		if filter.Year != 0 {
			query = query.Where("year = ?", filter.Year)
		}
		if filter.OwnerName != "" {
			query = query.Where("lower(owners.name) like ?", "%"+strings.ToLower(filter.OwnerName)+"%")
		}
		if filter.OwnerSurname != "" {
			query = query.Where("lower(owners.surname) like ?", "%"+strings.ToLower(filter.OwnerSurname)+"%")
		}
		if filter.OwnerPatronymic != "" {
			query = query.Where("lower(owners.patronymic) like ?", "%"+strings.ToLower(filter.OwnerPatronymic)+"%")
		}
	}

	err := query.Limit(limit).Offset(offset).Find(&cars).Error
	if err != nil {
		return nil, err
	}

	return cars, nil
}

// AddCarWithOwner adds a new car with its owner to the storage.
//
// Parameters:
// - car: a pointer to models.Cars representing the car to be added.
// - owner: a pointer to models.Owners representing the owner of the car.
// Returns an error if any.
func (s *Storage) AddCarWithOwner(car *models.Cars, owner *models.Owners) error {
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if owner exists
	var existingOwner *models.Owners
	if err := tx.Where("name = ? AND surname = ? AND patronymic = ?", owner.Name, owner.Surname, owner.Patronymic).First(&existingOwner).Error; err != nil {
		// owner does not exist, create one
		if err == gorm.ErrRecordNotFound {
			if err := tx.Create(owner).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	} else {
		owner = existingOwner
	}

	// create new car
	if err := tx.Create(&models.Cars{
		RegNum:  car.RegNum,
		Mark:    car.Mark,
		Model:   car.Model,
		Year:    car.Year,
		OwnerID: owner.ID,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// DeleteCarByRegNum deletes a car from the storage by its registration number.
//
// Parameters:
// - regNum: a string representing the registration number of the car to be deleted.
// Returns an error if any.
func (s *Storage) DeleteCarByRegNum(regNum string) error {
	car := &models.Cars{}
	s.Logger.Debug("Deleting car by regNum", zap.String("regNum", regNum))
	s.DB.Where("reg_num = ?", regNum).First(car)
	s.Logger.Debug("Car: ", zap.Any("car", car))

	if car.RegNum != regNum {
		return errors.New("car not found")
	}

	return s.DB.Delete(&car).Error
}

// UpdateCar updates a car in the storage with the information provided in the updatedCar parameter.
//
// Parameters:
// - updatedCar: a pointer to models.Cars representing the updated information of the car.
// Returns an error if any.
func (s *Storage) UpdateCar(updatedCar *models.Cars) error {
	var existingCar models.Cars

	/* if err := s.DB.First(&existingCar, updatedCar.RegNum).Error; err != nil {
		return err
	} */

	if err := s.DB.Where("reg_num = ?", updatedCar.RegNum).First(&existingCar).Error; err != nil {
		return err
	}

	return s.DB.Model(&existingCar).Updates(updatedCar).Error
}
