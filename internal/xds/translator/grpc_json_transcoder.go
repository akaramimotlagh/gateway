// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package translator

import (
	"errors"

	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	grpcjsontranscoder "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/grpc_json_transcoder/v3"
	hcmv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/envoyproxy/gateway/internal/ir"
	"github.com/envoyproxy/gateway/internal/xds/types"
)

func init() {
	registerHTTPFilter(&grpcJSONTranscoder{})
}

type grpcJSONTranscoder struct{}

var _ httpFilter = &grpcJSONTranscoder{}

// patchHCM builds and appends the gRPC-JSON transcoder filter to the HTTP Connection Manager if applicable.
func (*grpcJSONTranscoder) patchHCM(
	mgr *hcmv3.HttpConnectionManager,
	irListener *ir.HTTPListener,
) error {
	if mgr == nil {
		return errors.New("hcm is nil")
	}

	if irListener == nil {
		return errors.New("ir listener is nil")
	}

	if !listenerContainsGRPCJSONTranscoder(irListener) {
		return nil
	}

	// Return early if filter already exists.
	for _, httpFilter := range mgr.HttpFilters {
		if httpFilter.Name == wellknown.GRPCJSONTranscoder {
			return nil
		}
	}

	grpcJSONTranscoderFilter, err := buildHCMGRPCJSONTranscoderFilter()
	if err != nil {
		return err
	}

	mgr.HttpFilters = append([]*hcmv3.HttpFilter{grpcJSONTranscoderFilter}, mgr.HttpFilters...)

	return nil
}

// buildHCMGRPCJSONTranscoderFilter returns a gRPC-JSON transcoder filter.
func buildHCMGRPCJSONTranscoderFilter() (*hcmv3.HttpFilter, error) {
	grpcJSONTranscoderProto := &grpcjsontranscoder.GrpcJsonTranscoder{}

	grpcJSONTranscoderAny, err := anypb.New(grpcJSONTranscoderProto)
	if err != nil {
		return nil, err
	}

	return &hcmv3.HttpFilter{
		Name: wellknown.GRPCJSONTranscoder,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{
			TypedConfig: grpcJSONTranscoderAny,
		},
	}, nil
}

// listenerContainsGRPCJSONTranscoder returns true if the provided listener has gRPC-JSON transcoder
// policies attached to its routes.
func listenerContainsGRPCJSONTranscoder(irListener *ir.HTTPListener) bool {
	if irListener == nil {
		return false
	}

	for _, route := range irListener.Routes {
		if route.Traffic != nil && route.Traffic.GRPCJSONTranscoder != nil {
			return true
		}
	}

	return false
}

// patchRoute patches the provided route with the gRPC-JSON transcoder config if applicable.
func (*grpcJSONTranscoder) patchRoute(route *routev3.Route, irRoute *ir.HTTPRoute, _ *ir.HTTPListener) error {
	if route == nil {
		return errors.New("xds route is nil")
	}
	if irRoute == nil {
		return errors.New("ir route is nil")
	}
	if irRoute.Traffic == nil || irRoute.Traffic.GRPCJSONTranscoder == nil {
		return nil
	}

	filterCfg := route.GetTypedPerFilterConfig()
	if _, ok := filterCfg[wellknown.GRPCJSONTranscoder]; ok {
		// This should not happen since this is the only place where the gRPC-JSON transcoder
		// filter is added in a route.
		return errors.New("route already contains gRPC-JSON transcoder config")
	}

	transcoder := irRoute.Traffic.GRPCJSONTranscoder

	routeCfgProto := &grpcjsontranscoder.GrpcJsonTranscoder{
		DescriptorSet: &grpcjsontranscoder.GrpcJsonTranscoder_ProtoDescriptorBin{
			ProtoDescriptorBin: []byte(transcoder.ProtoDescriptor),
		},
	}

	// Configure services
	if len(transcoder.Services) > 0 {
		routeCfgProto.Services = transcoder.Services
	}

	// Configure print options
	if transcoder.PrintOptions != nil {
		routeCfgProto.PrintOptions = &grpcjsontranscoder.GrpcJsonTranscoder_PrintOptions{
			AddWhitespace:             transcoder.PrintOptions.AddWhitespace != nil && *transcoder.PrintOptions.AddWhitespace,
			AlwaysPrintPrimitiveFields: transcoder.PrintOptions.AlwaysPrintPrimitiveFields != nil && *transcoder.PrintOptions.AlwaysPrintPrimitiveFields,
			AlwaysPrintEnumsAsInts:    transcoder.PrintOptions.AlwaysPrintEnumsAsInts != nil && *transcoder.PrintOptions.AlwaysPrintEnumsAsInts,
			PreserveProtoFieldNames:   transcoder.PrintOptions.PreserveProtoFieldNames != nil && *transcoder.PrintOptions.PreserveProtoFieldNames,
		}
	}

	// Configure match incoming request route
	if transcoder.MatchIncomingRequestRoute != nil {
		routeCfgProto.MatchIncomingRequestRoute = *transcoder.MatchIncomingRequestRoute
	}

	// Configure ignored query parameters
	if len(transcoder.IgnoredQueryParameters) > 0 {
		routeCfgProto.IgnoredQueryParameters = transcoder.IgnoredQueryParameters
	}

	// Configure auto mapping
	if transcoder.AutoMapping != nil {
		routeCfgProto.AutoMapping = *transcoder.AutoMapping
	}

	// Configure ignore unknown query parameters
	if transcoder.IgnoreUnknownQueryParameters != nil {
		routeCfgProto.IgnoreUnknownQueryParameters = *transcoder.IgnoreUnknownQueryParameters
	}

	// Configure convert gRPC status
	if transcoder.ConvertGRPCStatus != nil {
		routeCfgProto.ConvertGrpcStatus = *transcoder.ConvertGRPCStatus
	}

	routeCfgAny, err := anypb.New(routeCfgProto)
	if err != nil {
		return err
	}

	if filterCfg == nil {
		route.TypedPerFilterConfig = make(map[string]*anypb.Any)
	}

	route.TypedPerFilterConfig[wellknown.GRPCJSONTranscoder] = routeCfgAny

	return nil
}

func (g *grpcJSONTranscoder) patchResources(*types.ResourceVersionTable, []*ir.HTTPRoute) error {
	return nil
}