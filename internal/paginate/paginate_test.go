package paginate

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestNewPaginateOptions(t *testing.T) {
	expected := QueryOptions{
		Limit:  25,
		Offset: 0,
		Order: Order{
			Field:  "id",
			Method: OrderMethodAsc,
		},
	}

	options := NewPaginateOptions()

	if !cmp.Equal(expected, options) {
		t.Fatal(cmp.Diff(expected, options))
	}
}

func TestWithOffset(t *testing.T) {
	t.Parallel()

	expected := NewPaginateOptions()
	expected.Offset = 15

	options := NewPaginateOptions(WithOffset(15))

	if !cmp.Equal(expected, options) {
		t.Fatal(cmp.Diff(expected, options))
	}
}

func TestWithLimit(t *testing.T) {
	t.Parallel()

	// Test valid
	expected := NewPaginateOptions()
	expected.Limit = 5
	options := NewPaginateOptions(WithLimit(5))

	if !cmp.Equal(expected, options) {
		t.Fatal(cmp.Diff(expected, options))
	}

	// Test invalid limit
	expected = NewPaginateOptions()
	options = NewPaginateOptions(WithLimit(0))

	if !cmp.Equal(expected, options) {
		t.Fatal(cmp.Diff(expected, options))
	}
}

func TestWithOrder(t *testing.T) {
	t.Parallel()

	// Test valid
	expected := NewPaginateOptions()
	expected.Order.Method = OrderMethodDesc
	options := NewPaginateOptions(WithOrder(OrderMethodDesc))

	if !cmp.Equal(expected, options) {
		t.Fatal(cmp.Diff(expected, options))
	}

	// Test invalid order
	expected = NewPaginateOptions()
	options = NewPaginateOptions(WithOrder("notAscOrDesc"))

	if !cmp.Equal(expected, options) {
		t.Fatal(cmp.Diff(expected, options))
	}
}

func TestWithOrderField(t *testing.T) {
	t.Parallel()

	expected := NewPaginateOptions()
	expected.Order.Field = "username"
	options := NewPaginateOptions(WithOrderField("username"))

	if !cmp.Equal(expected, options) {
		t.Fatal(cmp.Diff(expected, options))
	}
}
