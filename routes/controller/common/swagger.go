package common

import (
	"encoding/json"
	"teampilot/integrations/dba"
	"teampilot/structs/dts"
	"fmt"
	"os"
)

func InitializeRBAC() error {
	data, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		return err
	}

	var swagger map[string]interface{}
	if err := json.Unmarshal(data, &swagger); err != nil {
		return err
	}

	tx := dba.DB.Begin()

	paths := swagger["paths"].(map[string]interface{})
	for path, methods := range paths {
		for method, details := range methods.(map[string]interface{}) {
			details := details.(map[string]interface{})
			tags := details["tags"].([]interface{})

			perm := dts.Permission{
				Name:        fmt.Sprintf("%s %s", method, path),
				Path:        path,
				Method:      method,
				Description: fmt.Sprintf("Tags: %v", tags),
			}

			if err := tx.Create(&perm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}
