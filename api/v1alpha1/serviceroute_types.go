package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// ServiceEndpoint publishes a platform service that can be reached through a
// VPC attachment. The object describes the service contract; the controller
// that owns the service is responsible for its health and implementation.
//
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=svcendpoint
// +kubebuilder:printcolumn:name="ADDRESS",type="string",JSONPath=".spec.address"
// +kubebuilder:printcolumn:name="PORT",type="integer",JSONPath=".spec.port"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type ServiceEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceEndpointSpec `json:"spec,omitempty"`
}

// ServiceEndpointSpec defines a platform service address and its scope.
type ServiceEndpointSpec struct {
	// ServiceClass identifies the platform capability implemented by this
	// endpoint, for example "ipv6-egress-dns".
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	ServiceClass string `json:"serviceClass"`

	// Address is the stable address that consumers use to reach the service.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="isIP(self)",message="address must be an IP address"
	Address string `json:"address"`

	// Port is the service port.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// Protocol is the transport protocol accepted by the endpoint.
	// +kubebuilder:validation:Required
	Protocol NetworkRuleProtocol `json:"protocol"`

	// DeliveryMode defines where the service backing attachment must run. The
	// current dataplane supports only a backing attachment on every node that
	// has selected consumers; it does not forward private-service traffic to a
	// remote node.
	// +kubebuilder:validation:Required
	DeliveryMode ServiceEndpointDeliveryMode `json:"deliveryMode"`

	// AttachmentRef identifies the Cloud API VPCAttachment that hosts this
	// endpoint when the service is delivered through a private VPC path.
	// The network API keeps this reference opaque; the consuming controller
	// resolves it against the Cloud API.
	// +optional
	AttachmentRef *ServiceEndpointAttachmentReference `json:"attachmentRef,omitempty"`

	// Region limits endpoint selection to a region. An empty value means the
	// endpoint is not region-scoped.
	// +optional
	Region string `json:"region,omitempty"`
}

// ServiceEndpointDeliveryMode defines the placement contract for a service.
// +kubebuilder:validation:Enum=NodeLocal
type ServiceEndpointDeliveryMode string

const (
	// ServiceEndpointDeliveryModeNodeLocal requires the endpoint's backing
	// attachment to be present on each node where selected consumers run.
	ServiceEndpointDeliveryModeNodeLocal ServiceEndpointDeliveryMode = "NodeLocal"
)

// ServiceEndpointAttachmentReference identifies the private attachment that
// backs a platform service endpoint.
type ServiceEndpointAttachmentReference struct {
	// Namespace is the namespace containing the Cloud API VPCAttachment.
	// +kubebuilder:validation:MinLength=1
	Namespace string `json:"namespace"`

	// Name is the Cloud API VPCAttachment name.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// ServiceEndpointList is a list of ServiceEndpoint resources.
// +kubebuilder:object:root=true
type ServiceEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceEndpoint `json:"items"`
}

// ServiceRoutePolicy selects VPC attachments that should receive a route to a
// platform ServiceEndpoint. The selector is evaluated by the network control
// plane; Galactic compiles the selected attachments into dataplane state.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=svcroute
// +kubebuilder:printcolumn:name="SERVICE",type="string",JSONPath=".spec.serviceRef.name"
// +kubebuilder:printcolumn:name="REGION",type="string",JSONPath=".spec.region"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type ServiceRoutePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceRoutePolicySpec   `json:"spec,omitempty"`
	Status ServiceRoutePolicyStatus `json:"status,omitempty"`
}

// ServiceRoutePolicySpec defines service selection and attachment eligibility.
type ServiceRoutePolicySpec struct {
	// ServiceRef identifies the ServiceEndpoint in this namespace.
	// +kubebuilder:validation:Required
	ServiceRef ServiceEndpointReference `json:"serviceRef"`

	// AttachmentSelector selects eligible VPC attachments by authoritative
	// labels. Labels that grant access to platform services must be assigned by
	// the network control plane, not by tenants.
	// +kubebuilder:validation:Required
	AttachmentSelector metav1.LabelSelector `json:"attachmentSelector"`

	// ProtocolPorts limits the traffic that the route is intended to carry.
	// An empty list means the endpoint's declared port and protocol are used.
	// +optional
	// +listType=atomic
	ProtocolPorts []ServiceRouteProtocolPort `json:"protocolPorts,omitempty"`

	// Region limits selected attachments to a region. An empty value means the
	// policy applies in every region where the endpoint is available.
	// +optional
	Region string `json:"region,omitempty"`
}

// ServiceRouteProtocolPort identifies an allowed transport protocol and port.
type ServiceRouteProtocolPort struct {
	// Protocol is the transport protocol.
	// +kubebuilder:validation:Required
	Protocol NetworkRuleProtocol `json:"protocol"`

	// Port is the transport port.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`
}

// ServiceEndpointReference identifies a ServiceEndpoint in the policy's
// namespace.
type ServiceEndpointReference struct {
	// Name is the ServiceEndpoint name.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// ServiceRoutePolicyStatus contains durable policy conditions only. Per-
// attachment selection and programming state belongs in metrics and logs.
type ServiceRoutePolicyStatus struct {
	// ObservedGeneration is the .metadata.generation this status was computed from.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions contains durable acceptance and service-availability state.
	//
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

const (
	// ConditionTypeServiceRouteAccepted reports whether the policy is valid and
	// can be compiled by the dataplane controllers.
	ConditionTypeServiceRouteAccepted = "Accepted"
)

// ServiceRoutePolicyList is a list of ServiceRoutePolicy resources.
// +kubebuilder:object:root=true
type ServiceRoutePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceRoutePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&ServiceEndpoint{}, &ServiceEndpointList{},
		&ServiceRoutePolicy{}, &ServiceRoutePolicyList{},
	)
}
