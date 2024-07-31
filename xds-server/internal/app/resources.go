package app

import (
	"context"
	"errors"
	"fmt"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/types/known/anypb"

	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"

	logger "github.com/asishrs/proxyless-grpc-lb/common/pkg/logger"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"go.uber.org/zap"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/google/uuid"
)

type podEndPoint struct {
	IP   string
	Port int32
}

func getK8sEndPoints(serviceNames []string) (map[string][]podEndPoint, error) {
	k8sEndPoints := make(map[string][]podEndPoint)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	endPoints, err := clientset.CoreV1().Endpoints("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Logger.Error("Received error while trying to get EndPoints", zap.Error(err))
	}
	logger.Logger.Debug("Endpoint in the cluster", zap.Int("count", len(endPoints.Items)))
	for _, serviceName := range serviceNames {
		for _, endPoint := range endPoints.Items {
			name := endPoint.GetObjectMeta().GetName()
			if name == serviceName {
				var ips []string
				var ports []int32
				for _, subset := range endPoint.Subsets {
					for _, address := range subset.Addresses {
						ips = append(ips, address.IP)
					}
					for _, port := range subset.Ports {
						ports = append(ports, port.Port)
					}
				}
				logger.Logger.Debug("Endpoint", zap.String("name", name), zap.Any("IP Address", ips), zap.Any("Ports", ports))
				var podEndPoints []podEndPoint
				for _, port := range ports {
					for _, ip := range ips {
						podEndPoints = append(podEndPoints, podEndPoint{ip, port})
					}
				}
				k8sEndPoints[serviceName] = podEndPoints
			}
		}
	}
	return k8sEndPoints, nil
}

func clusterLoadAssignment(podEndPoints []podEndPoint, clusterName string, region string, zone string) []types.Resource {
	var lbs []*endpoint.LbEndpoint
	for _, podEndPoint := range podEndPoints {
		logger.Logger.Debug("Creating ENDPOINT", zap.String("host", podEndPoint.IP), zap.Int32("port", podEndPoint.Port))
		hst := &core.Address{Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Address:  podEndPoint.IP,
				Protocol: core.SocketAddress_TCP,
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: uint32(podEndPoint.Port),
				},
			},
		}}

		lbs = append(lbs, &endpoint.LbEndpoint{
			HostIdentifier: &endpoint.LbEndpoint_Endpoint{
				Endpoint: &endpoint.Endpoint{
					Address: hst,
				}},
			HealthStatus: core.HealthStatus_HEALTHY,
		})
	}

	eds := []types.Resource{
		&endpoint.ClusterLoadAssignment{
			ClusterName: clusterName,
			Endpoints: []*endpoint.LocalityLbEndpoints{{
				Locality: &core.Locality{
					Region: region,
					Zone:   zone,
				},
				Priority:            0,
				LoadBalancingWeight: &wrappers.UInt32Value{Value: uint32(1000)},
				LbEndpoints:         lbs,
			}},
		},
	}
	return eds
}

func createCluster(clusterName string) []types.Resource {
	logger.Logger.Debug("Creating CLUSTER", zap.String("name", clusterName))
	cls := []types.Resource{
		&cluster.Cluster{
			Name:                 clusterName,
			LbPolicy:             cluster.Cluster_ROUND_ROBIN,
			ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
			EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
				EdsConfig: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_Ads{},
				},
			},
		},
	}
	return cls
}

func createVirtualHost(virtualHostName, listenerName, clusterName string) *route.VirtualHost {
	logger.Logger.Debug("Creating RDS", zap.String("host name", virtualHostName))
	vh := &route.VirtualHost{
		Name:    virtualHostName,
		Domains: []string{listenerName},

		Routes: []*route.Route{{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: clusterName,
					},
				},
			},
		}}}
	return vh

}

func createRoute(routeConfigName, virtualHostName, listenerName, clusterName string) []types.Resource {
	vh := createVirtualHost(virtualHostName, listenerName, clusterName)
	rds := []types.Resource{
		&route.RouteConfiguration{
			Name:         routeConfigName,
			VirtualHosts: []*route.VirtualHost{vh},
		},
	}
	return rds
}
func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{
		ConfigSourceSpecifier: &core.ConfigSource_Ads{
			Ads: &core.AggregatedConfigSource{},
		}}

	return source
}

func createListener(listenerName string, clusterName string, routeConfigName string) []types.Resource {
	logger.Logger.Debug("Creating LISTENER", zap.String("name", listenerName))
	routerConfig, _ := anypb.New(&router.Router{})
	logger.Logger.Debug("RouterConfig", zap.String("router config", routerConfig.String()))
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: routeConfigName,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name:       "http-router",
			ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: routerConfig},
		}},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		panic(err)
	}

	lds := []types.Resource{
		&listener.Listener{
			Name: listenerName,
			ApiListener: &listener.ApiListener{
				ApiListener: pbst,
			},
			Address: &core.Address{
				Address: &core.Address_SocketAddress{
					SocketAddress: &core.SocketAddress{
						Protocol: core.SocketAddress_TCP,
						Address:  "0.0.0.0",
						PortSpecifier: &core.SocketAddress_PortValue{
							PortValue: 10000,
						},
					},
				},
			},
			FilterChains: []*listener.FilterChain{{
				Filters: []*listener.Filter{{
					Name: wellknown.HTTPConnectionManager,
					ConfigType: &listener.Filter_TypedConfig{
						TypedConfig: pbst,
					},
				}},
			}},
		}}
	return lds
}

// GenerateSnapshot creates snapshot for each service
func GenerateSnapshot(services []string) (*cache.Snapshot, error) {
	k8sEndPoints, err := getK8sEndPoints(services)
	if err != nil {
		logger.Logger.Error("Error while trying to get EndPoints from k8s cluster", zap.Error(err))
		return nil, errors.New("Error while trying to get EndPoints from k8s cluster")
	}

	logger.Logger.Debug("K8s", zap.Any("EndPoints", k8sEndPoints))

	var eds []types.Resource
	var cds []types.Resource
	var rds []types.Resource
	var lds []types.Resource
	for service, podEndPoints := range k8sEndPoints {
		logger.Logger.Debug("Creating new XDS Entry", zap.String("service", service))
		eds = append(eds, clusterLoadAssignment(podEndPoints, fmt.Sprintf("%s-cluster", service), "my-region", "my-zone")...)
		cds = append(cds, createCluster(fmt.Sprintf("%s-cluster", service))...)
		rds = append(rds, createRoute(fmt.Sprintf("%s-route", service), fmt.Sprintf("%s-vhost", service), fmt.Sprintf("%s-listener", service), fmt.Sprintf("%s-cluster", service))...)
		lds = append(lds, createListener(fmt.Sprintf("%s-listener", service), fmt.Sprintf("%s-cluster", service), fmt.Sprintf("%s-route", service))...)
	}

	version := uuid.New()
	logger.Logger.Debug("Creating Snapshot", zap.String("version", version.String()), zap.Any("EDS", eds), zap.Any("CDS", cds), zap.Any("RDS", rds), zap.Any("LDS", lds))
	resources := map[resource.Type][]types.Resource{
		resource.ClusterType:  cds,
		resource.RouteType:    rds,
		resource.ListenerType: lds,
		resource.EndpointType: eds,
	}

	snapshot, err := cache.NewSnapshot(version.String(), resources)
	if err != nil {
		logger.Logger.Error("new snapshot", zap.Error(err))
	}

	if err := snapshot.Consistent(); err != nil {
		logger.Logger.Error("Snapshot inconsistency", zap.Any("snapshot", snapshot), zap.Error(err))
	}
	return snapshot, nil
}
