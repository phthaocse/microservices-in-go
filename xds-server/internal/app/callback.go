package app

import (
	"context"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"sync"

	logger "github.com/asishrs/proxyless-grpc-lb/common/pkg/logger"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"go.uber.org/zap"
)

// Report type
func (cb *Callbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	logger.Logger.Debug("cb.Report()  callbacks", zap.Any("Fetches", cb.Fetches), zap.Any("Requests", cb.Requests))
}

// OnStreamOpen type
func (cb *Callbacks) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	logger.Logger.Debug("OnStreamOpen", zap.Int64("id", id), zap.String("type", typ))
	return nil
}

// OnStreamClosed type
func (cb *Callbacks) OnStreamClosed(id int64, node *core.Node) {
	logger.Logger.Debug("OnStreamClosed", zap.Int64("id", id))
}

// OnStreamRequest type
func (cb *Callbacks) OnStreamRequest(id int64, req *discovery.DiscoveryRequest) error {
	logger.Logger.Debug("OnStreamRequest", zap.Int64("id", id), zap.Any("Request", req))
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Requests++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	return nil
}

// OnFetchRequest type
func (cb *Callbacks) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	logger.Logger.Debug("OnFetchRequest", zap.Any("Request", req))
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Fetches++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	return nil
}
func (cb *Callbacks) OnStreamDeltaResponse(id int64, req *discovery.DeltaDiscoveryRequest, res *discovery.DeltaDiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.DeltaResponses++
}

// OnFetchResponse type
func (cb *Callbacks) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	logger.Logger.Debug("OnFetchResponse", zap.Any("Request", req), zap.Any("Response", resp))
}

func (cb *Callbacks) OnDeltaStreamClosed(id int64, node *core.Node) {
	logger.Logger.Debug("delta stream %d of node %s closed\n", zap.Any("id", id), zap.Any("node", node.Id))

}
func (cb *Callbacks) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {
	logger.Logger.Debug("delta stream %d open for %s\n", zap.Any("id", id), zap.Any("typ", typ))
	return nil
}

// OnStreamDeltaRequest invokes StreamDeltaResponseFunc
func (cb *Callbacks) OnStreamDeltaRequest(streamID int64, req *discovery.DeltaDiscoveryRequest) error {
	logger.Logger.Debug("delta stream request open for %s\n", zap.Any("stream id", streamID), zap.Any("req", req))
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.DeltaRequests++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}

	return nil
}

func (cb *Callbacks) OnStreamResponse(ctx context.Context, id int64, req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	logger.Logger.Debug("OnStreamResponse", zap.Int64("id", id), zap.Any("Request", req), zap.Any("Response ", resp))
	cb.Report()
}

// Callbacks for XD Server
type Callbacks struct {
	Signal         chan struct{}
	Fetches        int
	Requests       int
	DeltaRequests  int
	DeltaResponses int
	mu             sync.Mutex
}
