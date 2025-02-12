# tbls-driver-tailordb

[tbls driver](https://github.com/k1LoW/tbls#external-database-driver) for [TailorDB](https://docs.tailor.tech/guides/tailordb/overview) schema definition.

## Usage

```console
$ tbls doc tailordb://path/to/workspace.cue
```

```console
$ tbls doc tailordb://path/to/services/tailordb/tailordb.cue
```

## Example

| System | Database document |
| --- | --- |
| [CRM](https://github.com/tailor-platform/templates/tree/main/crm) | [example/crm](example/crm) |
| [Inventory Management System](https://github.com/tailor-platform/templates/tree/main/ims) | [example/ims](example/ims) |

## Support type

- [x] CUE
- [ ] Terraform

## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/tbls-driver-tailordb
```

**go install:**

```console
$ go install github.com/k1LoW/tbls-driver-tailordb@latest
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/tbls-driver-tailordb/releases)

