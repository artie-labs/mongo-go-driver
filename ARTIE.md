# Artie fork of mongo-go-driver

This is `mongodb/mongo-go-driver` at `v2.3.1` plus the BSON length checks from upstream `v2.9.2`. It exists because Artie can't move to `v2.9.2` yet.

## Why it exists

- On `v2.3.1`, a malformed reply from a MongoDB server panics inside the driver's connection pool goroutine, which has no `recover`. The whole process exits. This is one of the bugs tracked under CVE-2026-81521.
- Upstream fixed it in `v2.9.2`, but drivers from `v2.4.0` on refuse servers below MongoDB 4.2, and `v2.9.x` refuses servers below 4.4. Amazon DocumentDB 4.0 reports itself as MongoDB 4.0, so those drivers can't connect to it.
- `v2.3.1` is the last release that accepts DocumentDB 4.0, and upstream has no patched `v2.3.x`.

## What differs from upstream v2.3.1

All changes are in `x/bsonx/bsoncore` and match upstream `v2.9.2`:

| File | Change |
| --- | --- |
| `bsoncore.go` | `valueLength` rejects a length that wrapped negative |
| `document.go` | `Document.Validate` and `NewDocumentFromReader` reject a length below 5 |
| `array.go` | `Array.Validate` rejects a length below 5 |
| `artie_length_test.go` | Tests for the above; each input panics on plain `v2.3.1` |

Nothing else is changed, including the module path and the minimum server version.

## What is not backported

- **`Client.BulkWrite` namespace bug and the GridFS chunk deletion bug** (the other two bugs under CVE-2026-81521). Artie's services don't call `BulkWrite` or use GridFS. Don't start using either on this fork.
- Every other fix and feature between `v2.3.1` and `v2.9.2`.

## How to use it

```
replace go.mongodb.org/mongo-driver/v2 => github.com/artie-labs/mongo-go-driver/v2 <tag>
```

Branches follow the other Artie forks: `artie/v2.3.1` is the base, changes go in through pull requests against it, and releases are tagged `v2.3.1-artie.N` from that branch.

## When to delete it

When no Artie customer runs DocumentDB 4.0 (or MongoDB below 4.4). Then remove the `replace` lines, move to the current upstream driver, and archive this repo.
