{{- define "gvList" -}}
{{- $groupVersions := . -}}

# BGP API Reference

**API Group:** `{{ (index . 0).GroupVersionString }}`
**Stability:** Alpha — fields and defaults may change without deprecation notice

---

## 1. Overview

This project provides Kubernetes APIs for expressing BGP routing intent. It is an
API project: it defines resources, relationships, validation, and status
contracts. It does not define how routing intent is realized.

The fleet uses two BGP planes per node:

| Plane        | Purpose                                                            |
|--------------|--------------------------------------------------------------------|
| **Underlay** | IPv6 unicast fabric routing between nodes and top-of-rack switches |
| **Overlay**  | L2VPN EVPN distribution for tenant workloads                       |

Each plane is represented by a separate `BGPRouter` resource. `BGPPeer`,
`BGPAdvertisement`, and `BGPPolicy` resources target a router by direct
reference (`routerRef`) or label selector (`routerSelector`).

### API Group

| Group | Purpose |
|-------|---------|
| `{{ (index . 0).GroupVersionString }}` | BGP routing context, sessions, advertisements, and policies |

All resources are **Namespaced**.

---

## 2. CRD Reference

| Kind | Short Name | Targeting |
|------|------------|-----------|
| [BGPRouter](#bgprouter) | `bgpr` | — |
| [BGPPeer](#bgppeer) | `bgppr` | `routerRef` XOR `routerSelector` |
| [BGPAdvertisement](#bgpadvertisement) | `bgpadv` | `routerRef` only |
| [BGPPolicy](#bgppolicy) | `bgpp` | `routerRef` XOR `routerSelector` |
| [BGPVRFInstance](#bgpvrfinstance) | `bgpvrf` | `routerRef` XOR `routerSelector` |

---

## 3. Type Reference

{{ range $groupVersions }}
{{ template "gvDetails" . }}
{{ end }}

{{- end -}}
