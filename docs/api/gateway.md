# API Reference

## Packages
- [network.datumapis.com/v1alpha1](#networkdatumapiscomv1alpha1)


## network.datumapis.com/v1alpha1

Package v1alpha1 contains API Schema definitions for the network.datumapis.com/v1alpha1 API group.

### Resource Types
- [NetworkEgressPolicy](#networkegresspolicy)
- [NetworkGateway](#networkgateway)
- [NetworkRule](#networkrule)



#### NetworkEgressPolicy



NetworkEgressPolicy enables internet egress for a single tenant
VPC/VPCAttachment, served by the shared hyperconverged gateway engine's
masquerade (SNAT/PAT) datapath. Unlike NetworkRule, it carries no
VIP/backend/port: egress is on or off for a (vpcRef, vpcAttachmentRef)
pair, existence-implies-enabled, not a per-flow rule — because the
destination of an egress flow is an arbitrary internet address, not a
pre-configured backend list.

It is namespaced (deployed to galactic-system) and tenant-writable; like
NetworkRule, vpcRef/vpcAttachmentRef are opaque string identifiers because
the VPC API is owned by a separate companion operator, not this repo. An
admission webhook (implemented by the consuming controller) must verify
the requester is authorized for vpcRef/vpcAttachmentRef before a policy is
accepted — see the Accepted condition.

Presence of an accepted NetworkEgressPolicy resolves only *enablement*
(should this tenant reach the egress datapath at all) — a routing-layer
decision (does the tenant's VRF have a default route toward the shared
egress_sid locator), not a per-packet datapath lookup. *Isolation*
(preventing two tenants with colliding ULA source addresses from
colliding in the egress connection table) is a separate, datapath-level
concern resolved by tagging each flow with the tenant/VRF identifier
carried in the egress_sid locator's own Argument bits, not by anything in
this spec.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `NetworkEgressPolicy` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[NetworkEgressPolicySpec](#networkegresspolicyspec)_ |  |  |  |
| `status` _[NetworkEgressPolicyStatus](#networkegresspolicystatus)_ |  |  |  |


#### NetworkEgressPolicySpec



NetworkEgressPolicySpec defines the desired egress-enablement state for a
tenant VPC/VPCAttachment.



_Appears in:_
- [NetworkEgressPolicy](#networkegresspolicy)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `vpcRef` _string_ | VPCRef is the opaque identifier of the target VPC this policy applies<br />to. This repo does not own the VPC API and does not validate the<br />identifier beyond non-emptiness; the admission webhook of the<br />consuming controller is responsible for verifying the requester is<br />authorized for this VPC before the policy is accepted. |  | MinLength: 1 <br />Required: \{\} <br /> |
| `vpcAttachmentRef` _string_ | VPCAttachmentRef is the opaque identifier of the target<br />VPCAttachment this policy applies to. Like VPCRef, this is an opaque<br />string reference validated by the admission webhook, not by this API. |  | MinLength: 1 <br />Required: \{\} <br /> |


#### NetworkEgressPolicyStatus



NetworkEgressPolicyStatus defines the observed state of a
NetworkEgressPolicy.



_Appears in:_
- [NetworkEgressPolicy](#networkegresspolicy)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `observedGeneration` _integer_ | ObservedGeneration is the .metadata.generation this status was computed from. |  |  |
| `assignedGatewayNode` _string_ | AssignedGatewayNode is the name of the NetworkGateway-backed gateway<br />node this policy's tenant should route egress traffic through,<br />mirroring NetworkRule's own status.primaryNode field and computed<br />the same way: assigned_node = hash(vpcRef) % <gateway node count><br />(design plan §4.5 — a tenant's egress node and its primary ingress<br />node are the same node, by design, so both fields are computed by<br />the identical AssignPrimaryNode function). The controller consuming<br />this CRD sets this field exactly once, at creation.<br />This value must never be silently recomputed by a reconciler once<br />set, for the exact same reason NetworkRuleStatus.PrimaryNode's own<br />doc comment gives: recomputing it on a later reconcile can flip<br />which node a tenant's egress traffic routes through and cause an<br />avoidable traffic flap; a reconciler that observes a stale or<br />removed node here must surface that via a condition instead of<br />overwriting the value. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource,<br />including Accepted (see AcceptedReasonOwnershipVerified /<br />AcceptedReasonOwnershipDenied in rule_types.go, reused as-is here). |  |  |


#### NetworkGateway



NetworkGateway defines an XDP ingress NAT+LB gateway engine instance bound
to a single dedicated gateway-role node. Exactly one NetworkGateway exists
per gateway node (spec.targetRef.name is the Kubernetes node name),
mirroring the BGPRouter node-scoped root object pattern. NetworkRule
resources are assigned to a NetworkGateway via status.primaryNode.

There is no tunnel overlay in this design (an earlier Geneve-based
approach was superseded before this type shipped): the gateway's XDP
program does Full-NAT (DNAT the VIP to a backend Pod's address, SNAT the
client's source to status.sRv6Address) and pushes an SRv6 uSID outer
header addressed to the backend's worker node directly, so return traffic
(addressed to status.sRv6Address) arrives back at this same gateway node
over the ordinary SRv6 fabric — no compute-node encap agent, no tunnel
endpoint to publish. status.sRv6Address is advertised into BGP the same
way any workload prefix is (a BGPAdvertisement naming it, /128, Argument
0 — the value PR #740 reserves and forbids registering into any tenant
VRF, guaranteeing it never collides with a real tenant's Argument), so
every other node learns a real kernel SEG6 route to it for free through
the existing EVPN pipeline.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `NetworkGateway` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[NetworkGatewaySpec](#networkgatewayspec)_ |  |  |  |
| `status` _[NetworkGatewayStatus](#networkgatewaystatus)_ |  |  |  |


#### NetworkGatewaySpec



NetworkGatewaySpec defines the desired state of a NetworkGateway.



_Appears in:_
- [NetworkGateway](#networkgateway)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `targetRef` _[TargetRef](#targetref)_ | TargetRef identifies the Node this gateway engine executes on. |  | Required: \{\} <br /> |


#### NetworkGatewayStatus



NetworkGatewayStatus defines the observed state of a NetworkGateway.



_Appears in:_
- [NetworkGateway](#networkgateway)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `observedGeneration` _integer_ | ObservedGeneration is the .metadata.generation this status was computed from. |  |  |
| `sRv6Address` _string_ | SRv6Address is this gateway node's own SRv6-reachable IPv6 address,<br />used as the Full-NAT SNAT source for every ingress flow this node<br />translates. Backend Pods' replies are naturally routed back to it<br />over the ordinary SRv6 fabric (the same mechanism that routes any<br />other node's traffic), where this node's XDP program decapsulates<br />and un-NATs them using its own conn_table — there is no separate<br />tunnel endpoint or overlay device to publish. Populated by the<br />engine once it has computed the address (a uFMT 48+16 uSID over this<br />node's own BGPRouter locator/node-ID, at the reserved Argument 0)<br />and advertised it into BGP. |  |  |
| `egressAddress` _string_ | EgressAddress is this gateway node's own publicly-routable IPv6<br />address, used as the masquerade SNAT source for every egress flow<br />this node translates on behalf of tenant VPC backends reaching the<br />internet. Unlike SRv6Address (reachable only within the SRv6 fabric),<br />this address must additionally be reachable from the public internet<br />— an eBGP/uplink-peering concern outside this API. Operator-supplied<br />via GALACTIC_GATEWAY_EGRESS_ADDRESS; there is no in-cluster<br />derivation mechanism yet, the same gap SRv6Address itself has today.<br />A gateway node not offering egress leaves this field empty. |  |  |
| `egressSID` _string_ | EgressSID is this gateway node's own egress_sid uSID *locator*<br />(design plan §3.1) — the reserved Argument range's Block+Node-ID<br />portion tenant VRF default routes encapsulate toward. Unlike<br />EgressAddress (a plain, publicly-routable address, no uSID<br />structure), this is a real uSID: other nodes need a kernel route to<br />it before they can install a SEG6 encap route naming it as the<br />destination (the same reason SRv6Address is advertised into BGP),<br />so this is published and advertised the same way SRv6Address/<br />EgressAddress already are. Operator-supplied via<br />GALACTIC_GATEWAY_EGRESS_SID; a gateway node not offering egress<br />leaves this field empty, always paired with EgressAddress (both<br />set, or neither). |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource. |  |  |


#### NetworkRule



NetworkRule defines ingress load-balancing and NAT for a single tenant
VPC/VPCAttachment, served by the shared hyperconverged gateway engine.
It is namespaced (deployed to galactic-system) and tenant-writable; the
vpcRef/vpcAttachmentRef fields are opaque string identifiers because the
VPC API is owned by a separate companion operator, not this repo. An
admission webhook (implemented by the consuming controller) must verify
the requester is authorized for vpcRef/vpcAttachmentRef before a rule is
accepted — see the Accepted condition.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `NetworkRule` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[NetworkRuleSpec](#networkrulespec)_ |  |  |  |
| `status` _[NetworkRuleStatus](#networkrulestatus)_ |  |  |  |


#### NetworkRuleBackend



NetworkRuleBackend is a single backend endpoint that ingress traffic
matching a NetworkRule's VIP addresses is load-balanced to.



_Appears in:_
- [NetworkRuleSpec](#networkrulespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `address` _string_ | Address is the backend's IPv4 or IPv6 address. |  | MaxLength: 45 <br />Required: \{\} <br /> |
| `port` _integer_ | Port is the backend's destination port. |  | Maximum: 65535 <br />Minimum: 1 <br />Required: \{\} <br /> |


#### NetworkRuleProtocol

_Underlying type:_ _string_

NetworkRuleProtocol is the transport protocol matched by a NetworkRule's
ingress VIP.

_Validation:_
- Enum: [tcp udp]

_Appears in:_
- [NetworkRuleSpec](#networkrulespec)

| Field | Description |
| --- | --- |
| `tcp` | NetworkRuleProtocolTCP matches TCP traffic.<br /> |
| `udp` | NetworkRuleProtocolUDP matches UDP traffic.<br /> |


#### NetworkRuleSpec



NetworkRuleSpec defines the desired ingress load-balancing state for a
tenant VPC/VPCAttachment.



_Appears in:_
- [NetworkRule](#networkrule)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `vpcRef` _string_ | VPCRef is the opaque identifier of the target VPC this rule applies<br />to. This repo does not own the VPC API and does not validate the<br />identifier beyond non-emptiness; the admission webhook of the<br />consuming controller is responsible for verifying the requester is<br />authorized for this VPC before the rule is accepted. |  | MinLength: 1 <br />Required: \{\} <br /> |
| `vpcAttachmentRef` _string_ | VPCAttachmentRef is the opaque identifier of the target<br />VPCAttachment this rule applies to. Like VPCRef, this is an opaque<br />string reference validated by the admission webhook, not by this API. |  | MinLength: 1 <br />Required: \{\} <br /> |
| `vipAddresses` _string array_ | VIPAddresses is the list of ingress VIP addresses (IPv4 and/or IPv6)<br />this rule provisions on the assigned gateway node(s). |  | MaxItems: 8 <br />MinItems: 1 <br />Required: \{\} <br />items:MaxLength: 45 <br /> |
| `protocol` _[NetworkRuleProtocol](#networkruleprotocol)_ | Protocol is the transport protocol matched by VIPAddresses/Port. |  | Enum: [tcp udp] <br />Required: \{\} <br /> |
| `port` _integer_ | Port is the ingress port on VIPAddresses that this rule load-balances. |  | Maximum: 65535 <br />Minimum: 1 <br />Required: \{\} <br /> |
| `backends` _[NetworkRuleBackend](#networkrulebackend) array_ | Backends is the list of backend address:port targets that ingress<br />traffic matching VIPAddresses/Protocol/Port is load-balanced to. |  | MaxItems: 64 <br />MinItems: 1 <br />Required: \{\} <br /> |


#### NetworkRuleStatus



NetworkRuleStatus defines the observed state of a NetworkRule.



_Appears in:_
- [NetworkRule](#networkrule)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `observedGeneration` _integer_ | ObservedGeneration is the .metadata.generation this status was computed from. |  |  |
| `primaryNode` _string_ | PrimaryNode is the name of the NetworkGateway-backed gateway node<br />assigned to advertise this rule's VIPAddresses at the preferred BGP<br />local-preference, per the active-active model: primary_node =<br />hash(vpcRef) % <gateway node count>. The controller consuming this<br />CRD sets this field exactly once, at creation.<br />This value must never be silently recomputed by a reconciler once<br />set. Recomputing it on a later reconcile can flip which node is<br />primary for a live VIP and cause an avoidable traffic flap; a<br />reconciler that observes a stale or removed node here must surface<br />that via a condition instead of overwriting the value. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource. |  |  |


#### TargetRef



TargetRef identifies the execution target for a BGPRouter.
Supported values for kind: Node.



_Appears in:_
- [NetworkGatewaySpec](#networkgatewayspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `kind` _string_ | Kind is the target resource kind (e.g. Node). |  | MinLength: 1 <br /> |
| `name` _string_ | Name is the name of the target resource. |  | MinLength: 1 <br /> |


