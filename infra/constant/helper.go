package constant

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"payment_service/infra/log"
	"payment_service/models"
)

var MappedStatusByName map[string]int64
var MappedStatusById map[int64]string

func MapStatusFromDB(r *gorm.DB) {
	var statusList []models.Status

	// Run the query and load all rows into statusList
	result := r.Table("status").Find(&statusList)
	if result.Error != nil {

	}

	// Initialize the maps
	MappedStatusByName = make(map[string]int64)
	MappedStatusById = make(map[int64]string)

	// Fill the maps
	for _, s := range statusList {
		MappedStatusByName[s.Name] = s.ID
		MappedStatusById[s.ID] = s.Name
	}
}

func TranslateStatusByID(statusId int64) (string, error) {
	status, found := MappedStatusById[statusId]
	if !found {
		log.Logger.WithFields(logrus.Fields{
			"status": status,
		}).Errorf("TranslateStatusByID(statusId int64)")
		return "", errors.New("invalid status id")
	}

	return status, nil
}

func TranslateStatusByName(status string) (int64, error) {
	statusId, found := MappedStatusByName[status]
	if !found {
		log.Logger.WithFields(logrus.Fields{
			"status": status,
		}).Errorf("statusId, found := MappedStatusByName[status]")
		return 0, errors.New("invalid status name")
	}

	return statusId, nil
}
