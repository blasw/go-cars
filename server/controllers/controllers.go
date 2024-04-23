package controllers

import (
	"encoding/json"
	"go-cars/storage"
	"go-cars/storage/models"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type getCarsDto struct {
	Page            int    `form:"page"`
	Amount          int    `form:"amount"`
	RegNum          string `form:"regNum"`
	Mark            string `form:"mark"`
	Model           string `form:"model"`
	Year            int    `form:"year"`
	OwnerName       string `form:"name"`
	OwnerSurname    string `form:"surname"`
	OwnerPatronymic string `form:"patronymic"`
}

// @Summary		Get cars
// @Description	Retrieve cars based on optional filters
// @Tags			cars
// @Produce		json
// @Param			page		query		int		false	"Page number. Sets to 1 if not specified"
// @Param			amount		query		int		false	"Amount of cars per page. Sets to 10 if not specified"
// @Param			regNum		query		string	false	"Optional car's regNum filter"
// @Param			mark		query		string	false	"Optional car's mark filter"
// @Param			year		query		int		false	"Optional car's year filter"
// @Param			name		query		string	false	"Optional owner's name filter"
// @Param			surname		query		string	false	"Optional owner's surname filter"
// @Param			patronymic	query		string	false	"Optional owner's patronymic filter"
// @Success		200			{array}		models.Cars
// @Failure		400			{object}	map[string]string
// @Failure		500			{object}	map[string]string
// @Router			/get [get]
func GetCars(db *storage.Storage, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto getCarsDto
		if err := ctx.ShouldBindQuery(&dto); err != nil {
			logger.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// default values
		if dto.Page == 0 {
			dto.Page = 1
		}
		if dto.Amount == 0 {
			dto.Amount = 10
		}

		logger.Debug("GetCars requested", zap.Any("dto", dto))
		filters := &storage.CarFilter{
			RegNum:          dto.RegNum,
			Mark:            dto.Mark,
			Model:           dto.Model,
			Year:            dto.Year,
			OwnerName:       dto.OwnerName,
			OwnerSurname:    dto.OwnerSurname,
			OwnerPatronymic: dto.OwnerPatronymic,
		}

		offset := (dto.Page - 1) * dto.Amount

		cars, err := db.GetFilteredCars(filters, dto.Amount, offset)
		if err != nil {
			logger.Error("Unable to get filtered cars", zap.String("err: ", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, cars)
	}
}

type addCarsDto struct {
	RegNums []string `json:"regNums" binding:"required"`
}

type sourceServerResponse struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Owner  struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
	} `json:"owner"`
}

func getCarFromSourceServer(regNum string) (*sourceServerResponse, error) {
	resp, err := http.Get(os.Getenv("SRC_ADDR") + "/info?regNum=" + regNum)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var car sourceServerResponse

	jsonBytes, _ := io.ReadAll(resp.Body)

	if err = json.Unmarshal(jsonBytes, &car); err != nil {
		return nil, err
	}

	return &car, nil
}

// @Summary		Add cars
// @Description	Adds cars from source server by regNums
// @Tags			cars
// @Accept			json
// @Produce		json
// @Param			req	body	controllers.addCarsDto	true	"Cars' regNums"
// @Success		200
// @Failure		400	{object}	map[string]string
// @Failure		500	{object}	map[string]string
// @Router			/add [post]
func AddCars(db *storage.Storage, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto addCarsDto
		//validating request
		if err := ctx.ShouldBindJSON(&dto); err != nil {
			logger.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var errors []string
		var cars []*sourceServerResponse

		for _, regNum := range dto.RegNums {
			car, err := getCarFromSourceServer(regNum)
			if err != nil {
				logger.Error("Unable to get car from source server", zap.String("err: ", err.Error()))
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				cars = append(cars, car)
			}
		}

		for _, car := range cars {
			err := db.AddCarWithOwner(
				&models.Cars{
					RegNum: car.RegNum,
					Mark:   car.Mark,
					Model:  car.Model,
					Year:   car.Year,
				},
				&models.Owners{
					Name:       car.Owner.Name,
					Surname:    car.Owner.Surname,
					Patronymic: car.Owner.Patronymic,
				},
			)

			if err != nil {
				logger.Error("Unable to add car with owner", zap.String("err: ", err.Error()), zap.Any("error type: ", reflect.TypeOf(err)))
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					errors = append(errors, "Car with regNum "+car.RegNum+" already exists")
				} else {
					errors = append(errors, err.Error())
				}
			}
		}

		if len(errors) > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errors})
			return
		}

		ctx.Status(http.StatusOK)
	}
}

type deleteCarDto struct {
	RegNum string `form:"regNum" binding:"required"`
}

// @Summary		Delete car
// @Description	Deletes car selected by regNum
// @Tags			cars
// @Param			regNum	query	string	true	"Car's regNum"
// @Success		200
// @Failure		400	{object}	map[string]string
// @Router			/delete [delete]
func DeleteCarByRegNum(db *storage.Storage, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto deleteCarDto
		if err := ctx.ShouldBind(&dto); err != nil {
			logger.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.DeleteCarByRegNum(dto.RegNum)
		if err != nil {
			logger.Error("Unable to delete car by regNum", zap.String("err: ", err.Error()))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.Status(http.StatusOK)
	}
}

type updateCarDto struct {
	RegNum string `json:"regNum" binding:"required"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
}

// @Summary		Update car
// @Description	Updates provided fields of car selected by regNum
// @Tags			cars
// @Accept			json
// @Param			req	body	controllers.updateCarDto	true	"Updated information"
// @Success		200
// @Failure		400	{object}	map[string]string
// @Router			/update [patch]
func UpdateCarByRegNum(db *storage.Storage, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto updateCarDto
		if err := ctx.ShouldBindJSON(&dto); err != nil {
			logger.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.UpdateCar(&models.Cars{
			RegNum: dto.RegNum,
			Mark:   dto.Mark,
			Model:  dto.Model,
			Year:   dto.Year,
		})

		if err != nil {
			logger.Error("Unable to update car by regNum", zap.String("err: ", err.Error()))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.Status(http.StatusOK)
	}
}
