package paginate

import log "github.com/sirupsen/logrus"

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

type PaginateFuncOption func(options *PaginateQueryOptions)

// NewPaginateOptions returns a default set of paginate options.
func NewPaginateOptions(opts ...PaginateFuncOption) PaginateQueryOptions {
	paginateOptions := PaginateQueryOptions{
		Limit:  25,
		Offset: 0,
		Order: PaginateOrder{
			Field:  "id",
			Method: OrderMethodAsc,
		},
	}

	// Add options
	for _, opt := range opts {
		opt(&paginateOptions)
	}

	return paginateOptions
}

func WithLimit(limit uint) PaginateFuncOption {
	return func(p *PaginateQueryOptions) {
		if limit <= 0 || limit > 1000 {
			log.Warnf("limit set outside of bounds, retaining defaults: %v", limit)
			return
		}

		p.Limit = limit
	}
}

func WithOffset(offset uint) PaginateFuncOption {
	return func(p *PaginateQueryOptions) {
		p.Offset = offset
	}
}

func WithOrderField(orderField string) PaginateFuncOption {
	return func(p *PaginateQueryOptions) {
		p.Order.Field = orderField
	}
}

func WithOrder(method OrderMethod) PaginateFuncOption {
	return func(p *PaginateQueryOptions) {
		if method != OrderMethodAsc && method != OrderMethodDesc {
			log.Warnf("invalid order set, retaining default: %s", method)
			return
		}

		p.Order.Method = method
	}
}
