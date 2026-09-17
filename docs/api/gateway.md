# API Reference

## Packages
- [network.datumapis.com/v1alpha1](#networkdatumapiscomv1alpha1)


## network.datumapis.com/v1alpha1

Package v1alpha1 contains API Schema definitions for the network.datumapis.com/v1alpha1 API group.

### Resource Types
- [EgressShard](#egressshard)
- [NetworkEgressPolicy](#networkegresspolicy)
- [NetworkGateway](#networkgateway)
- [NetworkRule](#networkrule)
- [ServiceVIPBinding](#servicevipbinding)



#### EgressShard



EgressShard marks a single node as a member of the sharded, stateful egress
translation tier (galactic-nat) — a component deliberately kept off the
ingress load-balancer's own consistent-hash ring (see NetworkGateway):
tenant egress traffic (backend -> arbitrary internet destination) is a
different traffic pattern from ingress (fixed VIP, fixed backend pool)
and needs its own placement ring, own per-flow state, and its own
self-routing return path, entirely independent of any NetworkGateway node.

A shard serves one or both address families. NAT66 (IPv6 -> IPv6) and NAT64
(IPv6 -> IPv4, RFC 6146) are the same function — stateful egress PAT with a
VRF-scoped session table — over different families, so one shard object
describes both rather than there being a second, near-duplicate kind.
Spec.ShardAddressIPv6 and Spec.ShardAddressIPv4 are each assigned only for
the family this shard translates; a shard with no IPv4 address performs no
NAT64.

Every shard owns a dedicated, publicly-routable address per family it
serves, and a flow's allocated masquerade port lives within it — so a reply
is delivered to the correct shard by ordinary unicast routing alone, with
no hashing or cross-shard lookup on the return path at all (the "any node
can determine the owning shard from the tuple alone" property, satisfied by
construction rather than by a replicated hash table).

A shard is told which addresses to translate to; it does not choose them.
The controller that owns a cell claims one address per family from the
addressing service and writes it into this spec, which keeps the
addressing-service credential off every translating node and keeps the
allocation request out of the path that attaches a workload. Status reports
what the node's datapath is actually programmed with, so an unclaimed
address, a stale datapath, and a divergence between the two are each
distinguishable.

A shard names nothing that selects it. Which shards serve which consumer
intent is decided entirely by label selectors evaluated on the selecting
side (see the LabelEgressShard* keys), so no consumer-facing API appears in
this group and no reference points from a datapath resource back at the
resource that placed traffic on it.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `EgressShard` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[EgressShardSpec](#egressshardspec)_ |  |  |  |
| `status` _[EgressShardStatus](#egressshardstatus)_ |  |  |  |


#### EgressShardSpec



EgressShardSpec defines the desired state of an EgressShard.

The address fields are written by the controller that owns the cell, not by
the shard and not by a consumer. Each is optional: an address the addressing
service has not yet handed out is absent rather than blank-but-required, so
a shard object exists from the moment its node is labelled and gains its
identity afterwards.

Each is also write-once, and cannot be unassigned once assigned. The
datapath claims a reply by exact match against the address it translates to,
so reassigning one strands the return traffic of every flow already
established through it, with no drain and no dual-address grace period
available to cover the change. A shard holding the wrong address is deleted
and recreated instead, which breaks those flows at a moment someone chose.



_Appears in:_
- [EgressShard](#egressshard)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `targetRef` _[TargetRef](#targetref)_ | TargetRef identifies the Node this shard executes on. |  | Required: \{\} <br /> |
| `shardAddressIPv6` _string_ | ShardAddressIPv6 is the dedicated, publicly-routable IPv6 address this<br />shard translates to — every NAT66 masquerade port it allocates lives<br />within this address, so any node can route a reply to the owning shard<br />using ordinary unicast routing on it alone, with no per-flow state<br />lookup anywhere but that shard.<br />Empty means no IPv6 address is assigned to this shard. |  |  |
| `shardAddressIPv4` _string_ | ShardAddressIPv4 is the dedicated, publicly-routable IPv4 address this<br />shard translates to, and the source an IPv4-only destination sees.<br />Unlike ShardAddressIPv6, reachability for it is not established by a<br />BGPAdvertisement into the EVPN fabric: an IPv4 reply arrives from the<br />internet, so the underlay or an upstream announcement must attract this<br />address to this node.<br />Empty means no IPv4 address is assigned to this shard. |  |  |
| `nat64Prefix` _string_ | NAT64Prefix is the IPv6 prefix whose synthesized addresses this shard<br />translates to IPv4 — one Datum-operated Network-Specific Prefix, shared<br />fabric-wide, never per-tenant. It must be the prefix the resolver<br />synthesizes into; a shard translating for a different one is a<br />blackhole with no symptom on either side.<br />Set together with ShardAddressIPv4 or not at all: an address with no<br />prefix has nothing to translate for, and a prefix with no address has<br />nothing to translate into. Write-once for the same reason the addresses<br />are: a shard that starts translating a different prefix blackholes every<br />destination the resolver already synthesized into the old one. |  |  |


#### EgressShardStatus



EgressShardStatus defines the observed state of an EgressShard.

The address and prefix fields report what this node's datapath is programmed
with, not what it was asked for. A value here that disagrees with the spec is
a shard that has not converged; a value here with no counterpart in the spec
is a shard still translating to an address nothing assigns any more.



_Appears in:_
- [EgressShard](#egressshard)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `observedGeneration` _integer_ | ObservedGeneration is the .metadata.generation this status was computed from. |  |  |
| `shardSID` _string_ | ShardSID is this shard's own uSID locator — a real SRv6 uSID (unlike<br />the ShardAddress fields, which are plain routable addresses),<br />advertised into BGP the same way any other node-reachability route is<br />(a /128 BGPAdvertisement, no VRFID/Function) so every other node learns<br />a kernel SEG6 route toward it before installing a tenant VRF's egress<br />route against it. One SID serves both families: which translation a<br />packet gets is decided from the inner destination, not from a second<br />SID.<br />Still chosen by an operator and reported here rather than assigned in<br />spec, unlike the addresses: a value another node already uses silently<br />diverts that node's traffic, so the assignment belongs to the<br />addressing service, which does not hand out identifiers of this kind<br />yet. |  |  |
| `shardAddressIPv6` _string_ | ShardAddressIPv6 is the IPv6 masquerade source this shard's datapath is<br />programmed with. Empty means it translates no IPv6 flow. |  |  |
| `shardAddressIPv4` _string_ | ShardAddressIPv4 is the IPv4 masquerade source this shard's datapath is<br />programmed with. Empty means it performs no NAT64. Publishing it is also<br />what makes the underlay reachability prerequisite in<br />Spec.ShardAddressIPv4 checkable. |  |  |
| `nat64Prefix` _string_ | NAT64Prefix is the prefix this shard's datapath is programmed to<br />translate. Empty whenever ShardAddressIPv4 is empty. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource,<br />including Programmed (see ConditionTypeProgrammed). |  |  |




#### NetworkEgressPolicy



NetworkEgressPolicy enables internet egress for a single tenant
VPC/VPCAttachment, served by the sharded, stateful galactic-nat tier
(see EgressShard). Unlike NetworkRule, it carries no VIP/backend/port:
egress is on or off for a (vpcRef, vpcAttachmentRef) pair,
existence-implies-enabled, not a per-flow rule — because the destination
of an egress flow is an arbitrary internet address, not a pre-configured
backend list.

It is namespaced (deployed to galactic-system) and tenant-writable; like
NetworkRule, vpcRef/vpcAttachmentRef are opaque string identifiers because
the VPC API is owned by a separate companion operator, not this repo. An
admission webhook (implemented by the consuming controller) must verify
the requester is authorized for vpcRef/vpcAttachmentRef before a policy is
accepted — see the Accepted condition.

Presence of an accepted NetworkEgressPolicy resolves only *enablement*
(should this tenant's VRF get an egress route toward the shared
translation tier at all) — unlike this type's original design (superseded), there is
no single "assigned gateway node" to compute or pin: any EgressShard may
serve any tenant's flow, chosen by the shard-placement consistent-hash
ring (internal/maglev, keyed on (tenant VRFID, backend, destination) —
see EgressShard's doc comment), not by a per-tenant node assignment stored
here. *Isolation* (preventing two tenants with colliding ULA source
addresses from colliding in the egress connection table) is a separate,
datapath-level concern resolved by tagging each flow with the VRFID
carried in the tenant's own SRv6 Argument, not by anything in this spec.





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
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource,<br />including Accepted (see AcceptedReasonOwnershipVerified /<br />AcceptedReasonOwnershipDenied in rule_types.go, reused as-is here). |  |  |


#### NetworkGateway



NetworkGateway marks a single dedicated gateway-role node as running the
Maglev/DSR consistent-hash L4 load-balancer engine. Exactly one
NetworkGateway exists per gateway node (spec.targetRef.name is the
Kubernetes node name), mirroring the BGPRouter node-scoped root object
pattern. NetworkRule resources are served by every NetworkGateway in the
namespace equally (anycast — see NetworkRuleStatus's doc comment); this
object's only job is to identify which nodes participate at all and
surface each node's engine health via Conditions.

This design does no address rewriting on the load-balancing path at all
(DSR — Direct Server Return): the gateway's XDP program picks a backend
via consistent hashing on the flow's 5-tuple and pushes an SRv6 uSID outer
header addressed to the backend's worker node directly, untouched
otherwise. The backend node answers the client directly (see
ServiceVIPBinding) — reply traffic never re-enters this gateway node, so
unlike the Full-NAT design this type originally described, a gateway node
has no SNAT source address of its own to publish and nothing analogous to
sRv6Address/egressAddress/egressSID belongs on this status anymore.





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
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource. |  |  |


#### NetworkRule



NetworkRule defines ingress load-balancing for a single tenant
VPC/VPCAttachment, served by every NetworkGateway node identically
(anycast Direct Server Return — see NetworkGateway's doc comment). It is
namespaced (deployed to galactic-system) and tenant-writable; the
vpcRef/vpcAttachmentRef fields are opaque string identifiers because the
VPC API is owned by a separate companion operator, not this repo. An
admission webhook (implemented by the consuming controller) must verify
the requester is authorized for vpcRef/vpcAttachmentRef before a rule is
accepted — see the Accepted condition.

Unlike the earlier Full-NAT design this type originally described, there
is no primary/secondary gateway node for a rule: every NetworkGateway
advertises every accepted rule's vipAddresses at equal BGP preference,
consistent-hashes the same backend list to the same backend for the same
flow (internal/maglev), and forwards without rewriting anything —
backend selection never needs a single "owning" node the way Full-NAT's
SNAT-source model did.





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
- [ServiceVIPBindingSpec](#servicevipbindingspec)

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
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource. |  |  |


#### ServiceVIPBinding



ServiceVIPBinding drives one worker node's backend-side half of the
DSR/Maglev load-balancer datapath: it tells the node which service VIP a
specific local backend must be reachable on, so the backend can reply to
clients directly (the "Direct Server Return" this design depends on —
see NetworkGateway's doc comment). Written by the same controller that
already resolves a NetworkRule's backends to worker nodes/SRv6
information (galactic-gateway's usidresolver.go), one object per
(node, VIP, backend) triple; consumed by a per-node reconciler running
inside galactic-router's tenant role.

EgressKind decides which of two backend mechanisms this object drives,
mirroring the same veth/tap fork the SRv6 uSID decap datapath already
has (internal/plumbing/ebpf/usidmap's egress_kind field). Both mechanisms
now converge on the same VIP-boundary substitution
(BackendAddress:BackendPort for VIPAddress:Port at the SRv6 uSID TC-BPF
boundary, usid_ingress's inbound half / usid_egress's outbound half) —
required for both, not just tap, since a decapsulated ingress packet is
delivered into the owning tenant's own VRF routing table, which has no
route to an address bound outside that VRF (found live: see galactic's
ServiceVIPBindingReconciler doc comment):

  - veth (container backend): the node ALSO binds VIPAddress on its own
    galactic-vip0 dummy interface (internal/plumbing/vip's
    Bind/Unbind/Verify in galactic) — this alone does not deliver
    anything to the backend pod (galactic-vip0 lives in the node's root
    namespace, not the tenant's VRF), but still lets the node itself
    verifiably answer on the VIP.
  - tap (VM backend): there is no guest-side configuration capability in
    this repo by design (internal/cnitap's own doc comment), so the
    substitution above is this kind's *only* delivery mechanism — the
    guest OS never needs to know the VIP exists at all.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `ServiceVIPBinding` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ServiceVIPBindingSpec](#servicevipbindingspec)_ |  |  |  |
| `status` _[ServiceVIPBindingStatus](#servicevipbindingstatus)_ |  |  |  |


#### ServiceVIPBindingEgressKind

_Underlying type:_ _string_

ServiceVIPBindingEgressKind mirrors usidmap's EgressKindVeth/EgressKindTap
constants at the API layer — see ServiceVIPBinding's doc comment.

_Validation:_
- Enum: [veth tap]

_Appears in:_
- [ServiceVIPBindingSpec](#servicevipbindingspec)

| Field | Description |
| --- | --- |
| `veth` | ServiceVIPBindingEgressKindVeth selects the netns-bind mechanism.<br /> |
| `tap` | ServiceVIPBindingEgressKindTap selects the transparent tap-boundary<br />translation mechanism.<br /> |


#### ServiceVIPBindingSpec



ServiceVIPBindingSpec defines the desired VIP binding/translation state.



_Appears in:_
- [ServiceVIPBinding](#servicevipbinding)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `targetRef` _[TargetRef](#targetref)_ | TargetRef identifies the Node this binding applies to. |  | Required: \{\} <br /> |
| `vipAddress` _string_ | VIPAddress is the service VIP the backend must be reachable on. |  | Required: \{\} <br /> |
| `port` _integer_ | Port is the VIP-facing port traffic arrives on. |  | Maximum: 65535 <br />Minimum: 1 <br />Required: \{\} <br /> |
| `protocol` _[NetworkRuleProtocol](#networkruleprotocol)_ | Protocol is the transport protocol this binding applies to. |  | Enum: [tcp udp] <br />Required: \{\} <br /> |
| `backendAddress` _string_ | BackendAddress is the backend's own real address (its pod-netns<br />address for a veth backend, or its actual guest-facing address for a<br />tap backend) — the VIP-boundary substitution target for both kinds<br />now (see EgressKind's own doc comment for why veth needs this too,<br />not just tap). +optional at the API level for the same reason<br />EgressKind itself carries no matching CEL requirement; the<br />reconciler validates it's set for either kind before doing anything. |  |  |
| `backendPort` _integer_ | BackendPort is the backend's own real port, paired with<br />BackendAddress — required for both kinds, see that field's doc<br />comment. |  | Maximum: 65535 <br />Minimum: 1 <br /> |
| `egressKind` _[ServiceVIPBindingEgressKind](#servicevipbindingegresskind)_ | EgressKind selects which backend mechanism this binding drives. |  | Enum: [veth tap] <br />Required: \{\} <br /> |


#### ServiceVIPBindingStatus



ServiceVIPBindingStatus defines the observed state of a ServiceVIPBinding.



_Appears in:_
- [ServiceVIPBinding](#servicevipbinding)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `observedGeneration` _integer_ | ObservedGeneration is the .metadata.generation this status was computed from. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource,<br />including Bound (set True once internal/plumbing/vip.Verify or the<br />equivalent tap-translation-table check confirms the backend is<br />actually reachable on VIPAddress, not merely that the bind/table-write<br />call itself returned nil). |  |  |


#### TargetRef



TargetRef identifies the execution target for a BGPRouter.
Supported values for kind: Node.



_Appears in:_
- [EgressShardSpec](#egressshardspec)
- [NetworkGatewaySpec](#networkgatewayspec)
- [ServiceVIPBindingSpec](#servicevipbindingspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `kind` _string_ | Kind is the target resource kind (e.g. Node). |  | MinLength: 1 <br /> |
| `name` _string_ | Name is the name of the target resource. |  | MinLength: 1 <br /> |


