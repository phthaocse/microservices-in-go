package api

import (
	"context"
	"github.com/phthaocse/microservices-in-go/shipping/internal/application/core/domain"
)

type Application struct {
	//db ports.DBPort
}

func NewApplication() *Application {
	return &Application{
		//db: db,
	}
}

func (a Application) Create(ctx context.Context, shipping domain.Shipping) (domain.Shipping, error) {
	//err := a.db.Save(ctx, &payment)
	//if err != nil {
	//	return domain.Payment{}, err
	//}
	return domain.Shipping{}, nil
}
