package payment

import (
	"context"
	"github.com/huseyinbabal/microservices-proto/golang/payment"
	"github.com/huseyinbabal/microservices/order/internal/application/core/domain"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/xds" // To install the xds resolvers and balancers.
)

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	creds := insecure.NewCredentials()
	log.Info("connecting to payment service %s", paymentServiceUrl)
	conn, err := grpc.NewClient(paymentServiceUrl, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	log.Info("connected to payment service")
	client := payment.NewPaymentClient(conn)
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(ctx context.Context, order *domain.Order) error {
	_, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})
	return err
}
