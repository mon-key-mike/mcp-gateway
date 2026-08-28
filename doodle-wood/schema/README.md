# Product metadata schema

One record per product. Version-controlled. Feeds every surface.

```
schema/
├── product.schema.json     # the schema (JSON Schema draft 2020-12)
├── validate.py             # validates records + reports per-surface readiness
└── examples/
    └── campfire-monkey-wood-slice.json   # worked example against a real product
```

Run it:

```sh
pip install jsonschema
./validate.py
```

## The two surfaces, and which owns what

This is the part that matters most, because the two sites are built and run by
completely different toolchains and it is easy to let them drift.

| | **doodle-wood.com** | **shop.doodle-wood.com** |
|---|---|---|
| Role | Brand + story front door | Live operator e-commerce |
| Built/managed by | Lovable, via Claude Desktop cowork code | Shopify (operator-run) |
| Owns | narrative, origin story, source doodles, burn video | price, inventory, orders, fulfilment, tax |
| Has a cart? | **No — by design.** Its build spec explicitly forbids cart/checkout | Yes |
| Reads from this schema | `narrative.*`, `media.*`, `classification.character`, `surfaces.brand.*` | `commerce.*`, `variants[]`, `physical.weight_g`, `seo.*`, `surfaces.shop.*` |
| Status today | live | live on myshopify; custom domain **not connected** |

**The rule:** the brand site never quotes a price it computes itself. It renders
`commerce.price` from this record and links to `surfaces.brand.deep_link_to_shop`.
If a price is wrong on the brand site, the fix is in this record, not in Lovable.

**Canonical URLs matter here.** Both sites will describe the same product. Set
`seo.canonical_url` to the shop URL for anything transactional, or the two
properties compete for the same search term and split their own ranking.

## Why the schema is mostly optional fields

Only four fields are required: `id`, `sku`, `title`, `status`. Everything else can
be `null`. That is deliberate — a record should be creatable the moment a doodle
becomes a product idea, long before anyone knows its weight or COGS.

The cost of that permissiveness is that "valid" stops meaning "ready". Two things
carry the weight instead:

1. **`_gaps[]`** — an explicit list of what is knowingly missing. A `null` is
   ambiguous (unknown? not applicable?); naming it in `_gaps` is not.
2. **`validate.py` surface readiness** — per surface, what is still missing before
   the product can publish there. This is the real gate.

Current output for the example product:

```
OK    dw-campfire-monkey-wood-slice
        shop       NOT READY  missing: media.hero_image, media.alt_text, physical.weight_g
        brand      NOT READY  missing: media.hero_image, surfaces.brand.deep_link_to_shop
```

That is honest: the product is live in Shopify but is not actually ready, because
it has no photograph. The schema surfaces that instead of hiding it.

## `id` is the join key

`id` (`dw-campfire-monkey-wood-slice`) never changes and is never reused. It is what
lets a Notion page, a Shopify product GID, a Lovable component, a booth QR code and
an affiliate link all refer to the same thing. Shopify GIDs are stored *inside* the
record (`surfaces.shop.shopify_product_gid`) rather than used as the primary key,
so migrating platforms does not orphan the data.

## Where the fields come from

| Block | Owner / source of truth |
|---|---|
| `narrative`, `classification.character` | Notion Company Library, Doodle Concept Register |
| `commerce`, `variants`, `surfaces.shop` | Live Shopify (pull with `../scripts/snapshot.sh`) |
| `production`, `physical` | Studio — nobody else knows these; currently mostly blank |
| `channels.affiliate_links` | Partner agreements; empty and expected to stay so for now |
| `sources.*` | Provenance. Every claim should trace back to one of these |

## Adding a product

1. Copy `examples/campfire-monkey-wood-slice.json` to `../catalog/products/<id>.json`.
2. Fill what you know. Leave the rest `null`.
3. List what you left out in `_gaps`.
4. `./validate.py` — fix schema errors, ignore readiness warnings until you intend to launch it.
5. Commit. The diff is the changelog.

## Known limitation

Nothing writes this schema *into* Shopify yet. `snapshot.sh` pulls live state one
way (Shopify → repo). A push direction — mapping these records onto Shopify
metafields so `production.*` and `channels.*` live on the product itself — is the
obvious next build, and would let the operator store answer questions it currently
can't. Not built, deliberately: the fields are almost all blank, so there is nothing
worth pushing yet.
