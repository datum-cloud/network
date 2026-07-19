# API Reference

## Packages
- [network.datumapis.com/v1alpha1](#networkdatumapiscomv1alpha1)


## network.datumapis.com/v1alpha1

Package v1alpha1 contains API Schema definitions for the network.datumapis.com/v1alpha1 API group.

### Resource Types
- [BGPAdvertisement](#bgpadvertisement)
- [BGPPeer](#bgppeer)
- [BGPPolicy](#bgppolicy)
- [BGPRouter](#bgprouter)
- [BGPVRFInstance](#bgpvrfinstance)



#### AFI

_Underlying type:_ _string_

AFI is the Address Family Indicator for a BGP address family.

_Validation:_
- Enum: [ipv4 ipv6 l2vpn]

_Appears in:_
- [AddressFamily](#addressfamily)

| Field | Description |
| --- | --- |
| `ipv4` |  |
| `ipv6` |  |
| `l2vpn` |  |


#### ASPathFilter



ASPathFilter matches BGP routes by AS path using a regex pattern.



_Appears in:_
- [BGPPolicyMatch](#bgppolicymatch)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `pattern` _string_ | Pattern is a regular expression matched against the AS path.<br />The AS path is represented as a space-separated string of ASNs. |  | MaxLength: 256 <br />MinLength: 1 <br /> |
| `matchType` _[ASPathMatchType](#aspathmatchtype)_ | MatchType determines whether the pattern must match the full path<br />or can match a substring. Default: "contains". | contains | Enum: [full contains] <br /> |


#### ASPathMatchType

_Underlying type:_ _string_

ASPathMatchType determines how an AS path filter pattern is applied.

_Validation:_
- Enum: [full contains]

_Appears in:_
- [ASPathFilter](#aspathfilter)

| Field | Description |
| --- | --- |
| `full` | ASPathMatchFull requires the pattern to match the entire AS path.<br /> |
| `contains` | ASPathMatchContains requires the pattern to match a substring of the AS path.<br /> |


#### AddressFamily



AddressFamily is a BGP address family expressed as an AFI/SAFI pair.
Valid combinations: ipv4/unicast, ipv6/unicast, l2vpn/evpn.



_Appears in:_
- [BGPAdvertisementSpec](#bgpadvertisementspec)
- [BGPPeerSpec](#bgppeerspec)
- [BGPPolicyMatch](#bgppolicymatch)
- [BGPRouterSpec](#bgprouterspec)
- [ResolvedRouterConfig](#resolvedrouterconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `afi` _[AFI](#afi)_ | AFI is the address family indicator. |  | Enum: [ipv4 ipv6 l2vpn] <br /> |
| `safi` _[SAFI](#safi)_ | SAFI is the subsequent address family indicator. |  | Enum: [unicast evpn] <br /> |


#### AdvertisementOriginateFrom



AdvertisementOriginateFrom defines how routes are sourced from a local system
resource rather than from the static Prefixes list.



_Appears in:_
- [BGPAdvertisementSpec](#bgpadvertisementspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[AdvertisementOriginateType](#advertisementoriginatetype)_ | Type is the source from which routes are originated. |  | Enum: [interface kernel] <br /> |
| `interfaceName` _string_ | InterfaceName is the name of the interface to originate routes from.<br />Required when type is "interface". |  | MinLength: 1 <br /> |


#### AdvertisementOriginateType

_Underlying type:_ _string_

AdvertisementOriginateType is the source from which routes are originated.

_Validation:_
- Enum: [interface kernel]

_Appears in:_
- [AdvertisementOriginateFrom](#advertisementoriginatefrom)

| Field | Description |
| --- | --- |
| `interface` | OriginateTypeInterface originates routes from local interface addresses.<br /> |
| `kernel` | OriginateTypeKernel originates routes learned from the kernel routing table.<br /> |


#### AdvertisementPolicyRef



AdvertisementPolicyRef references a BGPPolicy by name to apply as a conditional
filter before advertisement.



_Appears in:_
- [BGPAdvertisementSpec](#bgpadvertisementspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name is the name of the BGPPolicy within the same namespace. |  | MinLength: 1 <br /> |


#### AsPathSet



AsPathSet defines AS path manipulation operations.



_Appears in:_
- [BGPPolicySetActions](#bgppolicysetactions)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `prepend` _integer_ | Prepend adds an ASN to the AS path N times. |  | Maximum: 10 <br />Minimum: 1 <br /> |
| `asn` _integer_ | ASN is the AS number to prepend (used when prepend is set).<br />Defaults to the local ASN if not specified. |  |  |
| `replace` _integer array_ | Replace replaces the entire AS path with the given list.<br />Mutually exclusive with prepend. |  | MaxItems: 32 <br /> |


#### BGPAdvertisement



BGPAdvertisement defines routing information to advertise from a single BGPRouter.
Routes are originated from static Prefixes, redistributed routing table entries,
or local interface/kernel routes. Fan-out via routerSelector is intentionally not
supported to avoid ambiguous prefix attribution across multiple routers.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `BGPAdvertisement` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BGPAdvertisementSpec](#bgpadvertisementspec)_ |  |  |  |
| `status` _[BGPAdvertisementStatus](#bgpadvertisementstatus)_ |  |  |  |


#### BGPAdvertisementSpec



BGPAdvertisementSpec defines the desired advertisement state.
At least one of prefixes, redistribute, or originateFrom must be specified.



_Appears in:_
- [BGPAdvertisement](#bgpadvertisement)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `routerRef` _[RouterRef](#routerref)_ | RouterRef targets a single BGPRouter by name. |  | Required: \{\} <br /> |
| `addressFamily` _[AddressFamily](#addressfamily)_ | AddressFamily defines the AFI/SAFI for this advertisement. |  | Required: \{\} <br /> |
| `prefixes` _[Prefix](#prefix) array_ | Prefixes is the list of CIDR prefixes to advertise.<br />At least one of Prefixes, Redistribute, or OriginateFrom must be specified. |  | MaxItems: 256 <br />MaxLength: 64 <br /> |
| `redistribute` _[RedistributeSource](#redistributesource) array_ | Redistribute defines routing table sources to redistribute into BGP.<br />Routes matching the source type are originated without requiring explicit CIDR entries.<br />At least one of Prefixes, Redistribute, or OriginateFrom must be specified. |  | Enum: [static connected kernel] <br />MaxItems: 3 <br /> |
| `originateFrom` _[AdvertisementOriginateFrom](#advertisementoriginatefrom)_ | OriginateFrom defines how routes are sourced from a local interface or kernel table.<br />When set, routes are originated from the specified source in addition to any static Prefixes.<br />At least one of Prefixes, Redistribute, or OriginateFrom must be specified. |  |  |
| `policyRef` _[AdvertisementPolicyRef](#advertisementpolicyref)_ | PolicyRef references a BGPPolicy to apply as a conditional filter before advertisement.<br />Only routes that match the policy are originated. |  |  |
| `communities` _[Community](#community) array_ | Communities is the default list of BGP communities to attach to all advertised prefixes.<br />Per-prefix communities in Prefixes[n].communities replace this value for individual prefixes. |  | MaxItems: 64 <br />MaxLength: 32 <br /> |
| `localPreference` _integer_ | LocalPreference sets the default BGP LOCAL_PREF attribute for all advertised prefixes.<br />Per-prefix localPreference in Prefixes[n].localPreference overrides this value.<br />Only meaningful for iBGP sessions. |  | Minimum: 0 <br /> |
| `vrfID` _integer_ | VRFID is the 16-bit PoP-local VRF identifier this advertisement's SRv6<br />SID resolves into, per RFC 9800 uSID Argument addressing. Required when<br />Function is set. |  | Maximum: 65535 <br />Minimum: 1 <br /> |
| `function` _[SRv6Function](#srv6function)_ | Function is the RFC 8986 SRv6 endpoint behavior applied when this<br />advertisement originates a compressed uSID. Required when VRFID is set. |  | Enum: [End.DT4 End.DT6 End.DT46] <br /> |


#### BGPAdvertisementStatus



BGPAdvertisementStatus defines the observed state of BGPAdvertisement.



_Appears in:_
- [BGPAdvertisement](#bgpadvertisement)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `observedGeneration` _integer_ | ObservedGeneration is the .metadata.generation this status was computed from. |  |  |
| `advertisedPrefixes` _integer_ | AdvertisedPrefixes is the count of prefixes currently being originated. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource. |  |  |


#### BGPOrigin

_Underlying type:_ _string_

BGPOrigin is the BGP origin attribute per RFC 4271.

_Validation:_
- Enum: [igp egp incomplete]

_Appears in:_
- [BGPPolicySetActions](#bgppolicysetactions)

| Field | Description |
| --- | --- |
| `igp` | BGPOriginIGP indicates the route was learned via an IGP.<br /> |
| `egp` | BGPOriginEGP indicates the route was learned via the EGP protocol.<br /> |
| `incomplete` | BGPOriginIncomplete indicates the route origin is unknown.<br /> |


#### BGPPeer



BGPPeer defines a BGP session to a remote peer. It binds to one or more
BGPRouter instances via routerRef or routerSelector.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `BGPPeer` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BGPPeerSpec](#bgppeerspec)_ |  |  |  |
| `status` _[BGPPeerStatus](#bgppeerstatus)_ |  |  |  |


#### BGPPeerAuthentication



BGPPeerAuthentication defines authentication configuration for a BGP peer.

Supports both secret-ref-based (MD5/TCP-AO) and plain-text passwords.
The plain-text password is primarily for non-production environments.



_Appears in:_
- [BGPPeerSpec](#bgppeerspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `password` _string_ | Password is a plain-text authentication password.<br />Prefer AuthSecretRef for production deployments. |  |  |


#### BGPPeerBFD



BGPPeerBFD defines BFD (Bidirectional Forwarding Detection) parameters for a BGP peer.

BFD provides sub-second failure detection for BGP sessions, enabling fast
convergence without relying on BGP hold timers (default 90s).



_Appears in:_
- [BGPPeerSpec](#bgppeerspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled indicates whether BFD is enabled for this peer. |  |  |
| `minimumTx` _integer_ | MinimumTx is the minimum transmission interval in microseconds.<br />FRR / GoBGP use microseconds — the controller converts to milliseconds<br />when programming the runtime. Defaults to 300000 (300ms). |  |  |
| `minimumRx` _integer_ | MinimumRx is the minimum reception interval in microseconds.<br />Defaults to 300000 (300ms). |  |  |
| `detectMultiplier` _integer_ | DetectMultiplier is the number of missed packets before declaring the<br />BFD session down. Default is 3. |  | Maximum: 50 <br />Minimum: 3 <br /> |
| `multiHop` _boolean_ | MultiHop enables BFD for eBGP multi-hop sessions (RFC 5883).<br />Must be true when ebgpMultiHop is true. |  |  |


#### BGPPeerGracefulRestart



BGPPeerGracefulRestart defines BGP graceful restart parameters for a peer.

Graceful restart (RFC 4724) allows a BGP speaker to restart its control plane
without tearing down sessions. Routes are preserved during the restart period,
preventing MAC/IP flapping in EVPN deployments.



_Appears in:_
- [BGPPeerSpec](#bgppeerspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled indicates whether graceful restart is enabled for this peer. |  |  |
| `restartTime` _integer_ | RestartTime is the maximum time (in seconds) the peer expects the<br />local system to complete a graceful restart. Range: 1-1200.<br />Default: 120 seconds. |  | Maximum: 1200 <br />Minimum: 1 <br /> |


#### BGPPeerSpec



BGPPeerSpec defines the desired state of BGPPeer.



_Appears in:_
- [BGPPeer](#bgppeer)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `routerRef` _[RouterRef](#routerref)_ | RouterRef targets a single BGPRouter by name.<br />Mutually exclusive with routerSelector. |  |  |
| `routerSelector` _[RouterSelector](#routerselector)_ | RouterSelector targets one or more BGPRouters by label.<br />Mutually exclusive with routerRef. |  |  |
| `peerASN` _integer_ | PeerASN is the remote AS number. |  | Minimum: 1 <br />Required: \{\} <br /> |
| `address` _string_ | Address is the remote peer's IPv4 or IPv6 address. |  | MinLength: 1 <br />Required: \{\} <br /> |
| `remotePort` _integer_ | RemotePort is the TCP port used to establish the BGP session with<br />this peer. Defaults to 179 (the IANA-assigned BGP port) if unset. |  | Maximum: 65535 <br />Minimum: 1 <br /> |
| `description` _string_ | Description is a human-readable label for this peer (e.g., "spine-1"). |  |  |
| `authSecretRef` _[LocalSecretRef](#localsecretref)_ | AuthSecretRef references a Secret in the same namespace containing the<br />MD5 TCP authentication password under the key "password". |  |  |
| `addressFamilies` _[AddressFamily](#addressfamily) array_ | AddressFamilies defines the address families negotiated on this session. |  | MinItems: 1 <br /> |
| `holdTime` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta)_ | HoldTime is the BGP hold timer. Must be 0 (disabled) or >= 3s.<br />Defaults to 90s if unset. |  |  |
| `keepaliveTime` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta)_ | KeepaliveTime is the BGP keepalive interval. Must be <= HoldTime / 3.<br />Defaults to 30s if unset. Maximum is 30s; the controller enforces<br />the keepalive <= holdTime / 3 constraint with an Accepted=False condition. |  |  |
| `ebgpMultiHop` _boolean_ | EbgpMultiHop enables eBGP multihop when the peer is not directly connected.<br />When true, the TTL for BGP packets is set to 64 (vs. 1 for direct). |  |  |
| `bfd` _[BGPPeerBFD](#bgppeerbfd)_ | BFD configures Bidirectional Forwarding Detection for fast failure<br />detection (sub-500ms). Without BFD, convergence relies on BGP hold<br />timers (default 90s). |  |  |
| `gracefulRestart` _[BGPPeerGracefulRestart](#bgppeergracefulrestart)_ | GracefulRestart configures BGP graceful restart for this peer.<br />Prevents route flapping during control-plane restarts. |  |  |
| `multiSession` _boolean_ | MultiSession enables BFD over MP-BGP (RFC 8935). When true, each<br />AFI/SAFI negotiates its own BFD session independently. When false<br />(default), a single BFD session covers all AFI/SAFIs. |  |  |
| `routeMapIn` _string_ | RouteMapIn is the name of a BGPPolicy term set applied to routes<br />received from this peer (import direction). The name must match<br />a BGPPolicy resource in the same namespace. |  |  |
| `routeMapOut` _string_ | RouteMapOut is the name of a BGPPolicy term set applied to routes<br />advertised to this peer (export direction). The name must match<br />a BGPPolicy resource in the same namespace. |  |  |
| `nextHopSelf` _boolean_ | NextHopSelf controls next-hop behavior for iBGP sessions.<br />When true, the local router's IP is used as next-hop for all<br />advertised routes. Common in EVPN iBGP peering. |  |  |
| `removePrivateAS` _integer_ | RemovePrivateAS strips private AS numbers (64512-65535, 4200000000-4294967294)<br />from the AS path on eBGP export. If set to a non-zero value, private ASNs<br />are replaced with the specified value before export. |  | Minimum: 0 <br /> |
| `defaultOriginRoute` _[OriginType](#origintype)_ | DefaultOriginRoute controls default route origination for this peer.<br />"igp" originates a default route with IGP origin.<br />"egp" originates a default route with EGP origin.<br />"incomplete" originates a default route with incomplete origin.<br />Empty or unset means no default route origination. |  | Enum: [igp egp incomplete] <br /> |
| `authentication` _[BGPPeerAuthentication](#bgppeerauthentication)_ | Authentication configures peer authentication.<br />AuthSecretRef takes precedence when both are set. |  |  |


#### BGPPeerState

_Underlying type:_ _string_

BGPPeerState represents the BGP Finite State Machine state of a session.



_Appears in:_
- [BGPPeerStatus](#bgppeerstatus)

| Field | Description |
| --- | --- |
| `Idle` |  |
| `Connect` |  |
| `Active` |  |
| `OpenSent` |  |
| `OpenConfirm` |  |
| `Established` |  |


#### BGPPeerStatus



BGPPeerStatus defines the observed state of BGPPeer.



_Appears in:_
- [BGPPeer](#bgppeer)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `observedGeneration` _integer_ | ObservedGeneration is the .metadata.generation this status was computed from. |  |  |
| `sessionState` _[BGPPeerState](#bgppeerstate)_ | SessionState is the current BGP FSM state of this session. |  |  |
| `lastEstablishedTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#time-v1-meta)_ | LastEstablishedTime is the timestamp of the most recent Established transition. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource. |  |  |


#### BGPPolicy



BGPPolicy defines composable, ordered routing policy statements applied to a
BGPRouter in a specific direction (import or export). It binds to one or more
BGPRouter instances via routerRef or routerSelector.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `BGPPolicy` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BGPPolicySpec](#bgppolicyspec)_ |  |  |  |
| `status` _[BGPPolicyStatus](#bgppolicystatus)_ |  |  |  |


#### BGPPolicyAction

_Underlying type:_ _string_

BGPPolicyAction is the disposition applied when a policy term matches.

_Validation:_
- Enum: [permit deny]

_Appears in:_
- [BGPPolicyTerm](#bgppolicyterm)

| Field | Description |
| --- | --- |
| `permit` | BGPPolicyActionPermit allows the route and optionally applies set actions.<br /> |
| `deny` | BGPPolicyActionDeny drops the route. Set actions must not be specified.<br /> |


#### BGPPolicyDirection

_Underlying type:_ _string_

BGPPolicyDirection is the direction in which a BGPPolicy is applied.

_Validation:_
- Enum: [import export]

_Appears in:_
- [BGPPolicySpec](#bgppolicyspec)

| Field | Description |
| --- | --- |
| `import` | BGPPolicyDirectionImport applies the policy to routes received from peers.<br /> |
| `export` | BGPPolicyDirectionExport applies the policy to routes advertised to peers.<br /> |


#### BGPPolicyMatch



BGPPolicyMatch defines the conditions under which a policy term fires.



_Appears in:_
- [BGPPolicyTerm](#bgppolicyterm)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `any` _boolean_ | Any matches all routes. When true, all other match fields are ignored. |  |  |
| `addressFamilies` _[AddressFamily](#addressfamily) array_ | AddressFamilies constrains the match to specific AFI/SAFI combinations.<br />If empty, all address families are matched. |  | MaxItems: 8 <br /> |
| `prefixList` _string array_ | PrefixList constrains the match to routes whose prefix matches one of<br />the given CIDR blocks. Each entry must be a valid IPv4 or IPv6 CIDR. |  | MaxItems: 256 <br />items:MaxLength: 43 <br /> |
| `asPathFilter` _[ASPathFilter](#aspathfilter)_ | ASPathFilter matches routes by AS path using a regex pattern.<br />The pattern is matched against the full AS path string (space-separated ASNs). |  |  |
| `communityMatch` _string array_ | CommunityMatch matches routes by BGP community.<br />Each entry is a community string in ASN:NN or IP:NN format. |  | MaxItems: 32 <br />items:MaxLength: 32 <br /> |
| `evpnRouteType` _[EVPNRouteType](#evpnroutetype) array_ | EVPNRouteType matches specific EVPN route types.<br />Only meaningful when l2vpn/evpn address family is configured. |  | Enum: [inclusiveMulticastEthernetTag macIPAdvertisement iPPrefixAdvertisement stickyMACAddress iPv6PrefixAdvertisement] <br />MaxItems: 5 <br /> |
| `vni` _integer_ | VNI matches routes by VNI (VXLAN Network Identifier).<br />Range: 0–16777215 (24-bit VNI). |  | Maximum: 1.6777215e+07 <br />Minimum: 0 <br /> |
| `macAddress` _string_ | MACAddress matches MAC-IP routes (EVPN Type-2) by MAC address.<br />Format: colon-separated hex bytes (e.g., "aa:bb:cc:dd:ee:ff"). |  | Pattern: `^([0-9a-fA-F]\{2\}:)\{5\}[0-9a-fA-F]\{2\}$` <br /> |
| `ipPrefix` _string_ | IPPrefix matches routes by exact IP prefix (CIDR notation). |  | MaxLength: 43 <br /> |
| `localPreference` _integer_ | LocalPreference matches routes by BGP LOCAL_PREF value. |  | Minimum: 0 <br /> |
| `med` _integer_ | MED matches routes by Multi-Exit Discriminator value. |  | Minimum: 0 <br /> |


#### BGPPolicySetActions



BGPPolicySetActions defines mutations applied when a term matches with action "permit".



_Appears in:_
- [BGPPolicyTerm](#bgppolicyterm)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `communities` _[CommunitySet](#communityset)_ | Communities defines community add/remove operations. |  |  |
| `localPreference` _integer_ | LocalPreference sets the LOCAL_PREF attribute.<br />Only meaningful on import (iBGP) or export to iBGP peers. |  | Minimum: 0 <br /> |
| `origin` _[BGPOrigin](#bgporigin)_ | Origin sets the BGP origin attribute. |  | Enum: [igp egp incomplete] <br /> |
| `asPath` _[AsPathSet](#aspathset)_ | AsPath manipulates the AS path (prepend or replace). |  |  |
| `nextHop` _[NextHopSet](#nexthopset)_ | NextHop overrides the next-hop attribute. |  |  |
| `extCommunities` _[ExtendedCommunitySet](#extendedcommunityset)_ | ExtCommunities defines extended community add/remove operations.<br />Each entry must be in a valid extended community format (ASN:NN, IP:NN,<br />or type-specific like "rt:65000:100"). |  |  |
| `metric` _integer_ | Metric sets the MED (Multi-Exit Discriminator) attribute. |  | Minimum: 0 <br /> |
| `color` _integer_ | Color sets the SRv6 policy color for path selection. |  | Minimum: 0 <br /> |


#### BGPPolicySpec



BGPPolicySpec defines the desired route policy state.



_Appears in:_
- [BGPPolicy](#bgppolicy)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `routerRef` _[RouterRef](#routerref)_ | RouterRef targets a single BGPRouter by name.<br />Mutually exclusive with routerSelector. |  |  |
| `routerSelector` _[RouterSelector](#routerselector)_ | RouterSelector targets one or more BGPRouters by label.<br />Mutually exclusive with routerRef. |  |  |
| `direction` _[BGPPolicyDirection](#bgppolicydirection)_ | Direction is the policy direction: import or export. |  | Enum: [import export] <br />Required: \{\} <br /> |
| `terms` _[BGPPolicyTerm](#bgppolicyterm) array_ | Terms is the ordered list of policy statements.<br />Evaluated from lowest to highest sequence number. |  | MaxItems: 32 <br />MinItems: 1 <br /> |


#### BGPPolicyStatus



BGPPolicyStatus defines the observed state of BGPPolicy.



_Appears in:_
- [BGPPolicy](#bgppolicy)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `observedGeneration` _integer_ | ObservedGeneration is the .metadata.generation this status was computed from. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource. |  |  |


#### BGPPolicyTerm



BGPPolicyTerm is a single ordered policy statement with match conditions and an action.



_Appears in:_
- [BGPPolicySpec](#bgppolicyspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `sequence` _integer_ | Sequence is the evaluation order. Lower values are evaluated first.<br />Must be unique within the policy. |  | Maximum: 65535 <br />Minimum: 1 <br /> |
| `match` _[BGPPolicyMatch](#bgppolicymatch)_ | Match defines the conditions under which this term fires. |  |  |
| `action` _[BGPPolicyAction](#bgppolicyaction)_ | Action is the disposition when this term matches. |  | Enum: [permit deny] <br /> |
| `set` _[BGPPolicySetActions](#bgppolicysetactions)_ | Set defines mutations applied when action is "permit".<br />Must not be set when action is "deny". |  |  |


#### BGPRouter



BGPRouter defines a logical BGP routing context. It abstracts a processing
instance bound to a specific execution context (e.g., a VRF or network
namespace) on a target node and acts as the primary ownership boundary for
BGPPeer, BGPAdvertisement, and BGPPolicy resources.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `BGPRouter` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BGPRouterSpec](#bgprouterspec)_ |  |  |  |
| `status` _[BGPRouterStatus](#bgprouterstatus)_ |  |  |  |


#### BGPRouterPeerSummary



BGPRouterPeerSummary summarizes BGP peer session counts.



_Appears in:_
- [BGPRouterStatus](#bgprouterstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `total` _integer_ | Total is the total number of configured peers. |  |  |
| `established` _integer_ | Established is the count of peers currently in the Established state. |  |  |


#### BGPRouterPhase

_Underlying type:_ _string_

BGPRouterPhase describes the lifecycle phase of a BGPRouter.



_Appears in:_
- [BGPRouterStatus](#bgprouterstatus)

| Field | Description |
| --- | --- |
| `Pending` |  |
| `Ready` |  |
| `Failed` |  |


#### BGPRouterSpec



BGPRouterSpec defines the desired state of a BGPRouter.



_Appears in:_
- [BGPRouter](#bgprouter)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `targetRef` _[TargetRef](#targetref)_ | TargetRef identifies the Node this router executes on. |  | Required: \{\} <br /> |
| `roles` _[RouterRole](#routerrole) array_ | Roles describes the functional roles this router participates in.<br />At least one role is required. |  | Enum: [fabric tenant transit] <br />MinItems: 1 <br /> |
| `localASN` _integer_ | LocalASN is the BGP Autonomous System Number for this router.<br />Must be a valid 2-byte or 4-byte ASN per RFC 6793. |  | Minimum: 1 <br />Required: \{\} <br /> |
| `routerID` _string_ | RouterID is a unique 32-bit identifier expressed in IPv4 dotted-decimal notation.<br />In an IPv6-only underlay this is a logical identifier only. |  | Format: ipv4 <br />Required: \{\} <br /> |
| `addressFamilies` _[AddressFamily](#addressfamily) array_ | AddressFamilies defines the address families this router activates. |  | MinItems: 1 <br /> |
| `srv6Locator` _string_ | SRv6Locator is the SRv6 locator block this router owns, expressed as an<br />IPv6 CIDR (e.g. "2001:db8:ff01::/48"). Individual SRv6 endpoint SIDs are<br />host addresses within this block. |  |  |
| `nodeID` _integer_ | NodeID is this router's 8-bit slot within its PoP's SRv6Locator block,<br />used for RFC 9800 NEXT-CSID compression. Unique within the PoP.<br />Values 0 and 255 are reserved. |  | Maximum: 254 <br />Minimum: 1 <br /> |


#### BGPRouterStatus



BGPRouterStatus defines the observed state of a BGPRouter.



_Appears in:_
- [BGPRouter](#bgprouter)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `phase` _[BGPRouterPhase](#bgprouterphase)_ | Phase is the high-level lifecycle state of this router. |  |  |
| `observedGeneration` _integer_ | ObservedGeneration is the .metadata.generation this status was computed from. |  |  |
| `roles` _[RouterRole](#routerrole) array_ | Roles reflects the active roles as observed by the implementation. |  | Enum: [fabric tenant transit] <br /> |
| `peers` _[BGPRouterPeerSummary](#bgprouterpeersummary)_ | Peers summarizes peer session counts. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions contains the standard conditions for this resource. |  |  |


#### BGPVRFInstance



BGPVRFInstance configures an L2VPN EVPN VRF on matched BGPRouters.
The referenced BGPRouter must have l2vpn-evpn in its addressFamilies.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `network.datumapis.com/v1alpha1` | | |
| `kind` _string_ | `BGPVRFInstance` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BGPVRFInstanceSpec](#bgpvrfinstancespec)_ |  |  |  |
| `status` _[BGPVRFInstanceStatus](#bgpvrfinstancestatus)_ |  |  |  |


#### BGPVRFInstanceSpec



BGPVRFInstanceSpec defines the desired VRF configuration.



_Appears in:_
- [BGPVRFInstance](#bgpvrfinstance)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `routerRef` _[RouterRef](#routerref)_ | RouterRef targets a single BGPRouter by name.<br />Mutually exclusive with routerSelector. |  |  |
| `routerSelector` _[RouterSelector](#routerselector)_ | RouterSelector targets one or more BGPRouters by label.<br />Mutually exclusive with routerRef. |  |  |
| `vrfID` _integer_ | VRFID is the 16-bit PoP-local VRF identifier used for RFC 9800 uSID<br />Argument addressing and to derive the RFC 4364 Type 1 Route<br />Distinguisher ("routerID:vrfID"). Unique per (VPC, PoP). Value 0 is<br />reserved. |  | Maximum: 65535 <br />Minimum: 1 <br />Required: \{\} <br /> |
| `importRouteTargets` _[RouteTarget](#routetarget) array_ | ImportRouteTargets is the list of BGP extended community route targets<br />used to import routes into this VRF. |  | MaxItems: 32 <br />MinItems: 1 <br /> |
| `exportRouteTargets` _[RouteTarget](#routetarget) array_ | ExportRouteTargets is the list of BGP extended community route targets<br />attached to routes exported from this VRF. |  | MaxItems: 32 <br />MinItems: 1 <br /> |


#### BGPVRFInstanceStatus



BGPVRFInstanceStatus defines the observed state of BGPVRFInstance.



_Appears in:_
- [BGPVRFInstance](#bgpvrfinstance)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions are top-level conditions for this BGPVRFInstance. |  |  |
| `routers` _[RouterStatus](#routerstatus) array_ | Routers holds per-router reconciliation status. |  |  |


#### Community

_Underlying type:_ _string_

Community is a BGP community in ASN:NN or IP:NN format.

_Validation:_
- MaxLength: 32

_Appears in:_
- [BGPAdvertisementSpec](#bgpadvertisementspec)



#### CommunitySet



CommunitySet defines community add and remove operations.



_Appears in:_
- [BGPPolicySetActions](#bgppolicysetactions)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `add` _string array_ | Add is a list of communities to attach. |  | MaxItems: 32 <br />items:MaxLength: 24 <br /> |
| `remove` _string array_ | Remove is a list of communities to strip. |  | MaxItems: 32 <br />items:MaxLength: 24 <br /> |


#### EVPNRouteType

_Underlying type:_ _string_

EVPNRouteType is an EVPN route type per RFC 7432.

_Validation:_
- Enum: [inclusiveMulticastEthernetTag macIPAdvertisement iPPrefixAdvertisement stickyMACAddress iPv6PrefixAdvertisement]

_Appears in:_
- [BGPPolicyMatch](#bgppolicymatch)

| Field | Description |
| --- | --- |
| `inclusiveMulticastEthernetTag` | EVPNRouteTypeInclusiveMulticastEthernetTag is Type-1: Inclusive Multicast Ethernet Tag route.<br /> |
| `macIPAdvertisement` | EVPNRouteTypeMACIPAdvertisement is Type-2: MAC-IP Advertisement route.<br /> |
| `iPPrefixAdvertisement` | EVPNRouteTypeIPPrefixAdvertisement is Type-3: IP Prefix Advertisement route.<br /> |
| `stickyMACAddress` | EVPNRouteTypeStickyMACAddress is Type-4: Sticky MAC Address route.<br /> |
| `iPv6PrefixAdvertisement` | EVPNRouteTypeIPv6PrefixAdvertisement is Type-5: IPv6 Prefix Advertisement route.<br /> |


#### ExtendedCommunitySet



ExtendedCommunitySet defines extended community add and remove operations.



_Appears in:_
- [BGPPolicySetActions](#bgppolicysetactions)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `add` _string array_ | Add is a list of extended communities to attach. |  | MaxItems: 32 <br />items:MaxLength: 64 <br /> |
| `remove` _string array_ | Remove is a list of extended communities to strip. |  | MaxItems: 32 <br />items:MaxLength: 64 <br /> |


#### LocalSecretRef



LocalSecretRef references a Secret within the same namespace.
Cross-namespace references are not supported.



_Appears in:_
- [BGPPeerSpec](#bgppeerspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name is the name of the Secret. |  | MinLength: 1 <br /> |


#### NextHopSet



NextHopSet defines next-hop attribute overrides.



_Appears in:_
- [BGPPolicySetActions](#bgppolicysetactions)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `self` _boolean_ | Self sets the next-hop to the local router's BGP peer address. |  |  |
| `address` _string_ | Address sets the next-hop to a specific IP address.<br />Mutually exclusive with self. |  | MaxLength: 45 <br /> |


#### OriginType

_Underlying type:_ _string_

OriginType defines the BGP origin attribute values.

_Validation:_
- Enum: [igp egp incomplete]

_Appears in:_
- [BGPPeerSpec](#bgppeerspec)

| Field | Description |
| --- | --- |
| `igp` |  |
| `egp` |  |
| `incomplete` |  |


#### Prefix

_Underlying type:_ _string_

Prefix is an IPv4 or IPv6 CIDR prefix.

_Validation:_
- MaxLength: 64

_Appears in:_
- [BGPAdvertisementSpec](#bgpadvertisementspec)



#### RedistributeSource

_Underlying type:_ _string_

RedistributeSource is a local routing table source to redistribute into BGP.

_Validation:_
- Enum: [static connected kernel]

_Appears in:_
- [BGPAdvertisementSpec](#bgpadvertisementspec)

| Field | Description |
| --- | --- |
| `static` | RedistributeSourceStatic redistributes statically configured routes.<br /> |
| `connected` | RedistributeSourceConnected redistributes directly connected interface routes.<br /> |
| `kernel` | RedistributeSourceKernel redistributes routes from the kernel routing table.<br /> |


#### ResolvedRouterConfig



ResolvedRouterConfig holds the configuration resolved and applied to a specific router.



_Appears in:_
- [RouterStatus](#routerstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `routerID` _string_ | RouterID resolved for this router. |  |  |
| `asNumber` _integer_ | ASNumber is the AS number configured. |  |  |
| `addressFamilies` _[AddressFamily](#addressfamily) array_ | AddressFamilies configured. |  |  |


#### RouteTarget



RouteTarget is a BGP extended community in "ASN:NN" or "IP:NN" format.



_Appears in:_
- [BGPVRFInstanceSpec](#bgpvrfinstancespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `value` _string_ | Value is the route target extended community string.<br />Format: "ASN:NN" (e.g. "65000:100") or "IP:NN" (e.g. "192.0.2.1:100"). |  | MaxLength: 21 <br />MinLength: 1 <br /> |


#### RouterRef



RouterRef is a direct reference to a single BGPRouter by name within the same namespace.



_Appears in:_
- [BGPAdvertisementSpec](#bgpadvertisementspec)
- [BGPPeerSpec](#bgppeerspec)
- [BGPPolicySpec](#bgppolicyspec)
- [BGPVRFInstanceSpec](#bgpvrfinstancespec)
- [RouterTarget](#routertarget)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name is the name of the BGPRouter. |  | MinLength: 1 <br /> |


#### RouterRole

_Underlying type:_ _string_

RouterRole defines the functional role of a BGPRouter within the network.

_Validation:_
- Enum: [fabric tenant transit]

_Appears in:_
- [BGPRouterSpec](#bgprouterspec)
- [BGPRouterStatus](#bgprouterstatus)

| Field | Description |
| --- | --- |
| `fabric` |  |
| `tenant` |  |
| `transit` |  |


#### RouterSelector



RouterSelector selects one or more BGPRouter resources by label within the same namespace.



_Appears in:_
- [BGPPeerSpec](#bgppeerspec)
- [BGPPolicySpec](#bgppolicyspec)
- [BGPVRFInstanceSpec](#bgpvrfinstancespec)
- [RouterTarget](#routertarget)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `matchLabels` _object (keys:string, values:string)_ | MatchLabels is a map of key/value label pairs to match. |  |  |
| `matchExpressions` _[LabelSelectorRequirement](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#labelselectorrequirement-v1-meta) array_ | MatchExpressions is a list of label selector requirements. |  |  |


#### RouterStatus



RouterStatus holds per-router reconciliation status used by BGPVRFInstance.



_Appears in:_
- [BGPVRFInstanceStatus](#bgpvrfinstancestatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `routerName` _string_ | RouterName is the name of the BGPRouter this entry describes. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | Conditions are the per-router conditions. |  |  |
| `resolvedConfig` _[ResolvedRouterConfig](#resolvedrouterconfig)_ | ResolvedConfig holds the configuration that was actually applied. |  |  |


#### RouterTarget



RouterTarget is embedded by resources that bind to one or more BGPRouters.
Exactly one of routerRef or routerSelector must be set.



_Appears in:_
- [BGPPeerSpec](#bgppeerspec)
- [BGPPolicySpec](#bgppolicyspec)
- [BGPVRFInstanceSpec](#bgpvrfinstancespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `routerRef` _[RouterRef](#routerref)_ | RouterRef targets a single BGPRouter by name.<br />Mutually exclusive with routerSelector. |  |  |
| `routerSelector` _[RouterSelector](#routerselector)_ | RouterSelector targets one or more BGPRouters by label.<br />Mutually exclusive with routerRef. |  |  |


#### SAFI

_Underlying type:_ _string_

SAFI is the Subsequent Address Family Indicator for a BGP address family.

_Validation:_
- Enum: [unicast evpn]

_Appears in:_
- [AddressFamily](#addressfamily)

| Field | Description |
| --- | --- |
| `unicast` |  |
| `evpn` |  |


#### SRv6Function

_Underlying type:_ _string_

SRv6Function is the RFC 8986 endpoint behavior applied to a decapsulated
SRv6 packet, addressed via the uSID Argument space per RFC 9800.

_Validation:_
- Enum: [End.DT4 End.DT6 End.DT46]

_Appears in:_
- [BGPAdvertisementSpec](#bgpadvertisementspec)

| Field | Description |
| --- | --- |
| `End.DT4` | SRv6FunctionEndDT4 decapsulates and looks up the packet in an IPv4 VRF.<br /> |
| `End.DT6` | SRv6FunctionEndDT6 decapsulates and looks up the packet in an IPv6 VRF.<br /> |
| `End.DT46` | SRv6FunctionEndDT46 decapsulates and looks up the packet in an IPv4 or<br />IPv6 VRF based on the inner packet's address family.<br /> |


#### TargetRef



TargetRef identifies the execution target for a BGPRouter.
Supported values for kind: Node.



_Appears in:_
- [BGPRouterSpec](#bgprouterspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `kind` _string_ | Kind is the target resource kind (e.g. Node). |  | MinLength: 1 <br /> |
| `name` _string_ | Name is the name of the target resource. |  | MinLength: 1 <br /> |


