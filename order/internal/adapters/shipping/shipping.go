package shipping

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	shipping ShippingServiceClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	creds := insecure.NewCredentials()
	log.Info("connecting to shipping service %s", shippingServiceUrl)
	conn, err := grpc.NewClient(shippingServiceUrl, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	log.Info("connected to shipping service")
	client := NewShippingServiceClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) CreateShipping(ctx context.Context, orderId int64) error {
	_, err := a.shipping.CreateShipping(ctx, &CreateShippingRequest{
		UserId:  1,
		OrderId: orderId,
	})
	return err
}
