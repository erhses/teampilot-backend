package dba

type (
	Sorter map[string]string

	PaginationInput struct {
		Limit  int    `json:"limit" example:"20"`
		Page   int    `json:"page" example:"0"`
		Sorter Sorter `json:"sorter"`
		IsAll  *bool  `json:"is_all"`
	}

	SorterInput struct {
		Sorter Sorter `json:"sorter"`
	}

	BaseFilter struct {
		DeletedAt []*string `json:"deleted_at"`
		CreatedAt []*string `json:"created_at"`
	}
	BaseModifierFilter struct {
		ModifierID *uint     `json:"modifier_id"`
		UpdatedAt  []*string `json:"updated_at"`
	}

	CursorInput struct {
		Limit      int  `form:"limit" json:"limit"`
		PreviousID uint `form:"previous_id" json:"previous_id"`
	}
)
