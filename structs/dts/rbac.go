package dts


type (
	Permission struct {
		Base
		Name        string `json:"name"`
		Path        string `json:"path"`
		Method      string `json:"method"`
		Description string `json:"description"`
	}
	Group struct {
		Base
		Key         string       `json:"key"`
		Name        string       `json:"name"`
		Description string       `json:"description"`
		Permissions []Permission `gorm:"many2many:group_permissions;" json:"permissions"`
	}
	Role struct {
		Base
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Groups      []Group `gorm:"many2many:role_groups;" json:"groups"`
	}
)