package dba

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Paginate is do pagination scope
func Paginate(input *PaginationInput) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if input != nil && input.IsAll != nil && *input.IsAll {
			for key, val := range input.Sorter {
				order := "ASC"
				if strings.HasPrefix(strings.ToLower(val), "desc") {
					order = "DESC"
				}
				db = db.Order(fmt.Sprintf("%s %s", key, order))
			}
			return db
		}

		tmpLimit := 20
		tmpOffset := 0
		if input != nil {
			tmpLimit = input.Limit
			tmpOffset = input.Page * input.Limit
			for key, val := range input.Sorter {
				order := "ASC"
				if strings.HasPrefix(strings.ToLower(val), "desc") {
					order = "DESC"
				}
				db = db.Order(fmt.Sprintf("%s %s", key, order))
			}
		}
		return db.Offset(tmpOffset).Limit(tmpLimit)
	}
}

// Sort is do sort scope
func Sort(input *Sorter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if input != nil {
			for key, val := range *input {
				order := "ASC"
				if strings.HasPrefix(strings.ToLower(val), "desc") {
					order = "DESC"
				}
				db = db.Order(fmt.Sprintf("%s %s", key, order))
			}
		}

		return db
	}
}

// Cursor pagination
func Cursor(input CursorInput) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id < ?", input.PreviousID).Order("id DESC").Limit(input.Limit)
	}
}
