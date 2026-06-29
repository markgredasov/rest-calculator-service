package domain

import "fmt"

var (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type Pagination struct {
	Page     *int
	PageSize *int
	Total    int
}

func NewPaginationUnitialized(page *int, pageSize *int) Pagination {
	p := Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    UnitializedTotal,
	}
	p.normalize()
	return p
}

func NewPagination(page *int, pageSize *int, total int) Pagination {
	p := Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
	p.normalize()
	return p
}

func (p *Pagination) normalize() {
	if p.Page == nil {
		p.Page = &defaultPage
	}

	if p.PageSize == nil {
		p.PageSize = &defaultPageSize
	}

	if *p.PageSize > maxPageSize {
		p.PageSize = &maxPageSize
	}
}

func (p *Pagination) Validate() error {
	if *p.Page < 1 {
		return fmt.Errorf("incorrect page value")
	}

	if *p.PageSize < 1 || *p.PageSize > 100 {
		return fmt.Errorf("incorrect pageSize value")
	}

	return nil
}

func (p *Pagination) GetLimitAndOffset() (int, int, error) {
	if p.Page == nil || p.PageSize == nil {
		return 0, 0, fmt.Errorf("invalid page or pageSize values")
	}

	limit := *p.PageSize
	offset := (*p.Page - 1) * *p.PageSize

	return limit, offset, nil
}
