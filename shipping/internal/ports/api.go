package ports

import (
	"context"
	"github.com/phthaocse/microservices-in-go/shipping/internal/application/core/domain"
)

type APIPort interface {
	Create(ctx context.Context, shipping domain.Shipping) (domain.Shipping, error)
}
