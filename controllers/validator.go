package controllers

import (
	"fmt"
	"golang-banking-management-system/models"

	//"net/http"
	"reflect"

	"gorm.io/gorm"
	//"github.com/gin-gonic/gin"
)

func validateStruct(s interface{}) error {
	v := reflect.ValueOf(s)

	// Handle pointer structs
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Loop through fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Check for zero value (empty string, 0, nil, etc.)
		if isZeroValue(field) {
			fieldName := v.Type().Field(i).Name
			return fmt.Errorf("insufficient data provided: missing %s", fieldName)
		}
	}
	return nil
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int64, reflect.Float64:
		return v.Interface() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}

func BankExists(db *gorm.DB, bankID string) bool {
	var count int64
	db.Model(&models.Bank{}).Where("pk_bank_id= ?", bankID).Count(&count)
	return count > 0
}

func CustomerExists(db *gorm.DB, customerId string) bool {
	var count int64
	db.Model(&models.User{}).Where("pk_customer_id= ?", customerId).Count(&count)
	return count > 0
}

func BankBranchExists(db *gorm.DB, branchId string) bool {
	var count int64
	db.Model(&models.BankBranch{}).Where("pk_bank_branch_id= ?", branchId).Count(&count)
	return count > 0
}
