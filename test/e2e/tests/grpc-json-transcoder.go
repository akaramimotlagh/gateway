// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

//go:build e2e

package tests

import (
	"testing"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/gateway-api/conformance/utils/http"
	"sigs.k8s.io/gateway-api/conformance/utils/kubernetes"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, EnvoyGatewayBackendTest)
}

var GrpcJsonTranscoderTest = suite.ConformanceTest{
	ShortName:   "GrpcJsonTranscoder",
	Description: "Uses the gRPC-JSON transcoder filter",
	Manifests:   []string{"testdata/grpc-json-transcoder.yaml"},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		t.Run("transcodes HTTP JSON", func(t *testing.T) {
			// The implementation of this function was copied verbatim from a different test.
			// TODO: Try to invoke EchoTwo via HTTP / JSON instead of doing what it currently does.
			ns := "gateway-conformance-infra"
			routeNN := types.NamespacedName{Name: "httproute-to-backend-fqdn-http2", Namespace: ns}
			gwNN := types.NamespacedName{Name: "same-namespace", Namespace: ns}
			gwAddr := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gwNN), routeNN)
			BackendMustBeAccepted(t, suite.Client, types.NamespacedName{Name: "backend-fqdn-http2", Namespace: ns})

			expectedResponse := http.ExpectedResponse{
				Request: http.Request{
					Path: "/backend-fqdn-http2",
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Namespace: ns,
			}

			http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, expectedResponse)
		})
	},
}
