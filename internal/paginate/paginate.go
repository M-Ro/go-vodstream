package paginate

import log "github.com/sirupsen/logrus"

type OrderMethod string

const (
	OrderMethodAsc  OrderMethod = "ASC"
	OrderMethodDesc             = "DESC"
)

type Order struct {
	Field  string
	Method OrderMethod
}

type QueryOptions struct {
	Limit  uint
	Offset uint
	Order  Order
}

type FuncOption func(options *QueryOptions)

// NewPaginateOptions returns a default set of paginate options.
func NewPaginateOptions(opts ...FuncOption) QueryOptions {
	paginateOptions := QueryOptions{
		Limit:  25,
		Offset: 0,
		Order: Order{
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

func WithLimit(limit uint) FuncOption {
	return func(p *QueryOptions) {
		if limit <= 0 || limit > 1000 {
			log.Warnf("limit set outside of bounds, retaining defaults: %v", limit)
			return
		}

		p.Limit = limit
	}
}

func WithOffset(offset uint) FuncOption {
	return func(p *QueryOptions) {
		p.Offset = offset
	}
}

func WithOrderField(orderField string) FuncOption {
	return func(p *QueryOptions) {
		p.Order.Field = orderField
	}
}

func WithOrder(method OrderMethod) FuncOption {
	return func(p *QueryOptions) {
		if method != OrderMethodAsc && method != OrderMethodDesc {
			log.Warnf("invalid order set, retaining default: %s", method)
			return
		}

		p.Order.Method = method
	}
}
