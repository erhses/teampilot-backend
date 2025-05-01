package common

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type PaginationResult struct {
	Total   int         `json:"total"`
	Summary interface{} `json:"summary"`
	Items   interface{} `json:"items"`
}

// type Paginate struct {
// 	Limit  int         `json:"limit"`
// 	Page   int         `json:"page"`
// 	Filter interface{} `json:"filter"`
// }

// func NewPaginate(limit int, page int) *Paginate {
// 	return &Paginate{limit: limit, page: page}
// }

// func (p *Paginate) PaginatedResult(db *gorm.DB) *gorm.DB {
// 	offset := (p.page - 1) * p.limit

// 	return db.Offset(offset).
// 		Limit(p.limit)
// }
