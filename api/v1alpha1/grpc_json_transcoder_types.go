// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v1alpha1

import gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

// GRPCJSONTranscoder defines the configuration for gRPC-JSON transcoding.
type GRPCJSONTranscoder struct {
	// ProtoDescriptor defines how to obtain the protocol buffer descriptor set.
	// This is required for the transcoder to understand the gRPC service definition.
	ProtoDescriptor ProtoDescriptor `json:"protoDescriptor"`

	// Services defines the gRPC services that should be transcoded.
	// If not specified, all services in the proto descriptor will be transcoded.
	// +optional
	Services []string `json:"services,omitempty"`

	// PrintOptions defines the output format options for JSON conversion.
	// +optional
	PrintOptions *JSONPrintOptions `json:"printOptions,omitempty"`

	// MatchIncomingRequestRoute enables matching the incoming request route pattern
	// against the gRPC service method. This is useful when the HTTP route pattern
	// doesn't exactly match the gRPC method pattern.
	// +optional
	MatchIncomingRequestRoute *bool `json:"matchIncomingRequestRoute,omitempty"`

	// IgnoredQueryParameters defines query parameters to ignore during transcoding.
	// +optional
	IgnoredQueryParameters []string `json:"ignoredQueryParameters,omitempty"`

	// AutoMapping enables automatic field mapping for HTTP query parameters and headers
	// to gRPC message fields when explicit field mapping is not present in the proto.
	// +optional
	AutoMapping *bool `json:"autoMapping,omitempty"`

	// IgnoreUnknownQueryParameters determines whether to ignore unknown query parameters.
	// If true, unknown query parameters are ignored; if false, the request is rejected.
	// +optional
	IgnoreUnknownQueryParameters *bool `json:"ignoreUnknownQueryParameters,omitempty"`

	// ConvertGRPCStatus enables converting gRPC status to HTTP status codes.
	// If true, gRPC status codes are converted to appropriate HTTP status codes.
	// +optional
	ConvertGRPCStatus *bool `json:"convertGrpcStatus,omitempty"`
}

// ProtoDescriptor defines how to obtain the protocol buffer descriptor set.
// +union
// +kubebuilder:validation:XValidation:rule="(self.type == 'Inline' && has(self.inline) && !has(self.valueRef)) || (self.type == 'ValueRef' && !has(self.inline) && has(self.valueRef))",message="Exactly one of inline or valueRef must be set with correct type."
type ProtoDescriptor struct {
	// Type is the type of method to use to read the proto descriptor.
	// Valid values are Inline and ValueRef, default is ValueRef.
	//
	// +kubebuilder:default=ValueRef
	// +kubebuilder:validation:Enum=Inline;ValueRef
	// +unionDiscriminator
	Type *ProtoDescriptorType `json:"type"`

	// Inline contains the proto descriptor as a base64-encoded string.
	// This should be the binary-encoded FileDescriptorSet.
	//
	// +optional
	Inline *string `json:"inline,omitempty"`

	// ValueRef is a reference to a local ConfigMap that contains the proto descriptor.
	// The value of key `proto-descriptor` in the ConfigMap will be used.
	// If the key is not found, the first value in the ConfigMap will be used.
	//
	// +optional
	ValueRef *gwapiv1.LocalObjectReference `json:"valueRef,omitempty"`
}

// ProtoDescriptorType specifies the type of proto descriptor source.
// +kubebuilder:validation:Enum=Inline;ValueRef
type ProtoDescriptorType string

const (
	// ProtoDescriptorTypeInline allows the user to specify the proto descriptor inline as base64.
	ProtoDescriptorTypeInline ProtoDescriptorType = "Inline"

	// ProtoDescriptorTypeValueRef allows the user to specify the proto descriptor via ConfigMap reference.
	ProtoDescriptorTypeValueRef ProtoDescriptorType = "ValueRef"
)

// JSONPrintOptions defines options for JSON output formatting.
type JSONPrintOptions struct {
	// AddWhitespace adds whitespace for pretty-printing JSON output.
	// +optional
	AddWhitespace *bool `json:"addWhitespace,omitempty"`

	// AlwaysPrintPrimitiveFields always prints primitive fields even if they have default values.
	// +optional
	AlwaysPrintPrimitiveFields *bool `json:"alwaysPrintPrimitiveFields,omitempty"`

	// AlwaysPrintEnumsAsInts always prints enum values as integers instead of strings.
	// +optional
	AlwaysPrintEnumsAsInts *bool `json:"alwaysPrintEnumsAsInts,omitempty"`

	// PreserveProtoFieldNames preserves proto field names in JSON output instead
	// of converting them to camelCase.
	// +optional
	PreserveProtoFieldNames *bool `json:"preserveProtoFieldNames,omitempty"`
}