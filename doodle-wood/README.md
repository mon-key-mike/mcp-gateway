# Doodle Wood — store as code

Version-controlled record of the Doodle Wood Shopify storefront, so the local
directory, Notion, and the live store stay in sync while the store is iterated on fast.

## What's here

```
doodle-wood/
├── STATE.md                  # verified current state + ranked blockers. Read this first.
├── catalog/
│   ├── products.json         # canonical catalog (hand-maintained, reviewed)
│   ├── collections.json      # collection structure + product mapping
│   └── live-*.json           # raw API snapshots, written by scripts/snapshot.sh
├── brand/storefront-copy.md  # homepage/about/social copy, mirrored from Notion
├── docs/DOMAIN-SETUP.md      # how to finish shop.doodle-wood.com (manual, blocked)
├── backup/                   # point-in-time memory snapshots
└── scripts/snapshot.sh       # pull live store state into catalog/
```

## The three sources of truth, and who wins

There are three copies of this company's state. They drift. The rule:

| Source | Authoritative for |
|---|---|
| **Live Shopify store** | products, prices, inventory, collections, domains |
| **Notion Company Library** | brand voice, positioning, strategy, pricing *policy* |
| **This directory** | the reconciliation — what was actually verified, and when |

When they disagree, the live store wins on facts and Notion wins on intent.
This directory records which is which so the disagreement is visible instead of silent.

> Note: Notion's own workbook declares itself sole source of truth, but the Doodle Wood
> Company Library lives outside it and is richer and more current. That conflict is
> flagged on the venture page in Notion and is not resolved here.

## The iteration loop

Run this every time the store changes:

```sh
export SHOPIFY_STORE=000yqx-de.myshopify.com
export SHOPIFY_ADMIN_TOKEN=shpat_...        # never commit
./scripts/snapshot.sh
git diff doodle-wood/catalog/               # exactly what changed in the store
git add doodle-wood && git commit -m "chore(doodle-wood): store snapshot $(date -u +%F)"
```

`snapshot.sh` exits loudly if any product still has no images, because that is the
current launch blocker.

## Secrets

No tokens in this directory, ever. `SHOPIFY_ADMIN_TOKEN` comes from the environment.
`.gitignore` blocks `.env` and `*token*` files under this tree.
