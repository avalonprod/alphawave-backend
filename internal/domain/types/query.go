package types

type PaginationQuery struct {
	Skip  int64
	Limit int64
}

func (p *PaginationQuery) GetLimit() *int64 {
	if p.Limit == 0 {
		return nil
	}
	return &p.Limit
}

func (p *PaginationQuery) GetSkip() *int64 {
	if p.Skip == 0 {
		return nil
	}
	return &p.Skip
}

type GetUsersByQuery struct {
	PaginationQuery
}

type GetMembersByQuery struct {
	PaginationQuery
}
