package grpc

import (
	"context"
	"github.com/phthaocse/microservices-in-go/shipping/internal/adapters/grpc/proto"
	log "github.com/sirupsen/logrus"
)

func (a Adapter) CreateShipping(ctx context.Context, in *proto.CreateShippingRequest) (*proto.CreateShippingResponse, error) {
	log.WithContext(ctx).Info("Creating shipping...")
	return &proto.CreateShippingResponse{Message: "success"}, nil
}
