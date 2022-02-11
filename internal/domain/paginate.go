package domain

type OrderMethod string

const (
	OrderMethodAsc  = "ASC"
	OrderMethodDesc = "DESC"
)

type PaginateOrder struct {
	Field  string
	Method OrderMethod
}

type PaginateQueryOptions struct {
	Limit  uint
	Offset uint
	Order  PaginateOrder
}

// NewPaginateOptions returns a default set of paginate options.
func NewPaginateOptions() PaginateQueryOptions {
	return PaginateQueryOptions{
		Limit:  25,
		Offset: 0,
		Order: PaginateOrder{
			Field:  "id",
			Method: OrderMethodAsc,
		},
	}
}
