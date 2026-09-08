# Doodle Wood store — state of record

**Last verified against the live store: 2026-09-08**
Everything below was read from the Shopify Admin API, not assumed.

## Store identity

| Field | Value | OK? |
|---|---|---|
| Shop name | `Doodle-Wood` | ok — but brand docs write it **Doodle Wood**, no hyphen |
| myshopify domain | `000yqx-de.myshopify.com` | working |
| **Primary domain** | **`shop.doodle-wood.com`, SSL enabled** | **done** |
| Plan | Basic | ok |
| Currency | USD | ok |
| Contact email | `metamonkeymike@gmail.com` | real, but not a doodle-wood.com address |
| Billing address | North Ridgeville, US | set |
| Ships to | US only | ok |
| Password protection | **disabled — store is public on the brand domain** | see blocker 2 |
| Live theme | Horizon | ok |

Fixed since 2026-08-28: custom domain connected, store renamed off the Shopify
default, contact email and billing address set.

## Catalog

**11 products, 4 collections.**

8 physical products, all ACTIVE, all with **zero images**:

| Collection | Products | Price range |
|---|---|---|
| Hand-Burned Wood Slices | 2 | $28–38 |
| Pendants & Wearables | 2 | $9–18 |
| Signs & Wall Art | 1 | $45 |
| Stickers & Camp Goods | 3 | $4–22 |

3 revenue-line products added 2026-09-08, all **DRAFT** (see "Why drafts" below):

| Product | SKU | Price | Notion's description |
|---|---|---|---|
| Burn Your Own — Workshop Seat | `DW-WS-SEAT` | $49 | "the highest-throughput item" |
| Commission a Piece — Deposit | `DW-CM-DEP` | $50 *(placeholder)* | deposit-gated pipeline |
| QR Plaque — Business | `DW-QR-B2B` | $72 | "the highest-value repeat B2B item" |

**0 orders. 0 customers.** Ever.

## Blockers, ranked

1. **Every physical product has zero images.** Unchanged since 2026-08-27. For a brand
   whose proposition is *how the burned piece looks*, this is the conversion mechanism,
   not decoration. Single biggest blocker.
2. **The store is public on `shop.doodle-wood.com` with no images and no password.**
   Connecting the domain did not cause this, but it converted the exposure from an
   obscure myshopify URL nobody had into your real brand domain. Anything linking
   from `doodle-wood.com`, any QR code, any bio link lands a customer here.
   Either set a storefront password until photography lands, or ship photos.
3. **Payment provider not confirmed.** `supportedDigitalWallets` is still empty, which
   on a US store suggests no provider is activated. Not confirmable via API — check
   Settings → Payments. 0 orders is consistent with both "no traffic" and "checkout
   is broken", and those need opposite responses.
4. **Free-shipping thresholds lose money.** Economy $6.69 free ≥ $39, Standard $8.00
   free ≥ $70, Express $15.00 — Shopify starter defaults, never set for this catalog.
   Worked arithmetic in `docs/FUNNEL.md`. Blocked on product weights, which are null
   across all 8 records (`physical.weight_g`). A scale and ten minutes unblocks it.
   *Last verified 2026-08-28; not re-checked 2026-09-08.*
5. **Policies written but publication unconfirmed.** Drafts in `policies/`. This
   connection cannot read policy state, so I cannot tell whether they were pasted.
   Two placeholders they were waiting on — contact email and business address — are
   now filled in the store and can be substituted in.

## Unresolved conflicts

Recorded rather than silently resolved, because each has two credible sources.

1. **Workshop price: $49 or $30?** The `on-brand` skill states $49. Notion and
   `docs/FUNNEL.md` say $30. Built at **$49** as the newer brand authority.
   `docs/FUNNEL.md` still says $30 and is wrong until this is settled.
2. **Wood species: driftwood or basswood/pine?** The `on-brand` skill says Doodle Wood
   is pyrography on **Lake Erie driftwood**, sealed with beeswax from the family's own
   hives. The 8 existing product records say **basswood** and **pine**. The three new
   products are written to the driftwood story. Either the older listings are stale or
   the brand doc overstates — someone who knows the bench needs to say which.
3. **Store name hyphenation.** Live store is `Doodle-Wood`; brand docs use
   `Doodle Wood`. Cosmetic, appears in checkout and transactional email.

## Why the three new products are drafts

Each is missing one fact that would make it dishonest to sell today:

- **Workshop seat** — no session dates exist. Publishing takes money for an
  unscheduled event.
- **Commission deposit** — the $50 is a **placeholder**. Notion records the pipeline
  as deposit-gated but never states the amount. Set the real figure before publishing.
- **QR plaque** — $72 is the correct floor per Notion, but B2B orders are quoted per
  batch, so the intake step needs to exist before the buy button does.

All three flip to ACTIVE in one call once those are answered.

## Relationship to doodle-wood.com

`doodle-wood.com` is a brand + story + email-capture site with no cart or checkout, by
its own build spec. This store is the commerce surface it never had. The apex `A`
record serves the live brand site — do not touch it.

Open item carried from Notion: `auth.users` on the brand site is empty, so no admin
exists and the captured email list is unreadable in-app.
