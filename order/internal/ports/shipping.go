package ports

import (
	"context"
)

type ShippingPort interface {
	CreateShipping(ctx context.Context, orderId int64) error
}
