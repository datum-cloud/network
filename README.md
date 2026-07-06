# Datum Network

Kubernetes CRDs for BGP routing intent. This project defines the API contract —
CRDs, validation rules, and status conditions — but ships no controller.
External implementations consume these types and realize the routing logic.

**API group:** `network.datumapis.com/v1alpha1`
**Stability:** Alpha

---

## Resources

| Kind               | Short name | Scope                            | Description                                                                                        |
|--------------------|------------|----------------------------------|----------------------------------------------------------------------------------------------------|
| `BGPRouter`        | `bgpr`     | —                                | BGP routing context: AS number, router ID, address families, node binding. One per plane per node. |
| `BGPPeer`          | `bgppr`    | `routerRef` XOR `routerSelector` | BGP session to a remote peer.                                                                      |
| `BGPAdvertisement` | `bgpadv`   | `routerRef` only                 | Prefixes to originate from a router.                                                               |
| `BGPPolicy`        | `bgpp`     | `routerRef` XOR `routerSelector` | Ordered import/export route filtering.                                                             |
| `BGPVRFInstance`   | `bgpvrf`   | `routerRef` XOR `routerSelector` | L2VPN EVPN VRF: route distinguisher, import/export route targets.                                  |

`BGPRouter` is the ownership root — all other resources bind to it via
`routerRef` (single router) or `routerSelector` (label-based multi-router).

---

## Quick start

### Install the CRDs

```bash
kubectl apply -k config/crd
```

Verify:

```bash
kubectl get crds | grep network.datumapis.com
```

### Example: underlay router and peer

```yaml
apiVersion: network.datumapis.com/v1alpha1
kind: BGPRouter
metadata:
  name: node-1-underlay
  namespace: default
spec:
  targetRef:
    kind: Node
    name: node-1
  roles:
    - fabric
  localASN: 65000
  routerID: "10.0.0.1"
  addressFamilies:
    - afi: ipv6
      safi: unicast
---
apiVersion: network.datumapis.com/v1alpha1
kind: BGPPeer
metadata:
  name: node-1-to-tor
  namespace: default
spec:
  routerRef:
    name: node-1-underlay
  address: "2001:db8:fabric::1"
  peerASN: 65000
  addressFamilies:
    - afi: ipv6
      safi: unicast
```

```bash
kubectl apply -f - <<EOF
# paste the YAML above
EOF
```

For a complete walkthrough see the [Getting Started guide](docs/getting-started.md).

---

## Requirements

- Kubernetes 1.28+ — CEL validation functions `isIP()` and `isCIDR()` are used in CRD schemas.

---

## Development

Install dev tools first:

```bash
task install
```

| Command                       | Description                           |
|-------------------------------|---------------------------------------|
| `task build`                  | `go fmt` + `go vet` + `go build`      |
| `task test:unit`              | Unit tests with coverage              |
| `task test:e2e`               | Kind cluster + Chainsaw e2e suite     |
| `task lint` / `task lint-fix` | golangci-lint + yamlfmt               |
| `task generate`               | Regenerate `zz_generated.deepcopy.go` |
| `task manifests`              | Regenerate CRD YAML in `config/crd/`  |
| `task ci`                     | Full pipeline: build → lint → test    |

---

## Documentation

- [Getting started](docs/getting-started.md) — install and create your first resources
- [BGP API reference](docs/api/bgp.md) — full CRD field definitions, conditions, validation rules
- [Enhancements](docs/enhancements/) — design proposals

---

## License

[AGPL-3.0](LICENSE) — source-available, copyleft.
