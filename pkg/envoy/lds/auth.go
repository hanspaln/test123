package lds

import (
	"fmt"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	xds_ext_authz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	xds_hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"

	"github.com/openservicemesh/osm/pkg/auth"
	"github.com/openservicemesh/osm/pkg/errcode"
)

func (lb *listenerBuilder) getExtAuthConfig() *auth.ExtAuthConfig {
	extAuthConfig := lb.cfg.GetInboundExternalAuthConfig()
	if extAuthConfig.Enable {
		return &extAuthConfig
	}
	return nil
}

// getExtAuthzHTTPFilter returns an envoy HttpFilter given an ExternAuthConfig configuration
func getExtAuthzHTTPFilter(extAuthConfig *auth.ExtAuthConfig) *xds_hcm.HttpFilter {
	extAuth := &xds_ext_authz.ExtAuthz{
		Services: &xds_ext_authz.ExtAuthz_GrpcService{
			GrpcService: &envoy_config_core_v3.GrpcService{
				TargetSpecifier: &envoy_config_core_v3.GrpcService_GoogleGrpc_{
					GoogleGrpc: &envoy_config_core_v3.GrpcService_GoogleGrpc{
						TargetUri: fmt.Sprintf("%s:%d",
							extAuthConfig.Address,
							extAuthConfig.Port),
						StatPrefix: extAuthConfig.StatPrefix,
					},
				},
				Timeout: ptypes.DurationProto(extAuthConfig.AuthzTimeout),
			},
		},
		TransportApiVersion: envoy_config_core_v3.ApiVersion_V3,
		WithRequestBody: &xds_ext_authz.BufferSettings{
			MaxRequestBytes:     8192,
			AllowPartialMessage: true,
		},
		FailureModeAllow: extAuthConfig.FailureModeAllow,
	}

	extAuthMarshalled, err := ptypes.MarshalAny(extAuth)
	if err != nil {
		log.Error().Err(err).Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrMarshallingXDSResource)).
			Msg("Failed to marshal External Authorization config")
	}

	return &xds_hcm.HttpFilter{
		Name: wellknown.HTTPExternalAuthorization,
		ConfigType: &xds_hcm.HttpFilter_TypedConfig{
			TypedConfig: extAuthMarshalled,
		},
	}
}
