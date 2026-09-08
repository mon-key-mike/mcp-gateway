# Conversion funnel — analysis and fix order

Verified against the live store 2026-08-28. Every number here was pulled from the
Shopify Admin API, not assumed.

> **Update 2026-09-08.** The three missing revenue lines below now exist as DRAFT
> products — see `../STATE.md`. Two corrections to this document:
>
> - **The workshop price here ($30) is disputed.** The `on-brand` skill states **$49**.
>   The product was built at $49. Settle this and update both places.
> - The domain and store-name blockers in the fix table are **done**.
>
> The shipping arithmetic below was last verified 2026-08-28 and has not been
> re-checked.

## The headline

**The store cannot sell its three best products, because they do not exist in it.**

Notion's own revenue record for Doodle Wood names:

| Revenue line | Notion's own description | Product in store? |
|---|---|---|
| $30 "Burn Your Own" workshop seats | **"the highest-throughput item"** | **No** |
| Commissions, deposit-gated | a named revenue path with a defined pipeline | **No** |
| B2B QR plaques from $72 | **"the highest-value repeat B2B item"** | **No** |

The store sells eight physical merch items, priced $4–$45. Every revenue line the
business describes as its best is unpurchasable. No amount of CTA tuning fixes a
missing product.

This is the answer to "what is the highest-leverage funnel work": it is not copy,
and it is not photography. It is that three revenue lines have no checkout path.

## Shipping economics — a live margin leak

Current rates (flat, Domestic zone, pulled live):

| Method | Rate | Free over |
|---|---|---|
| Economy | $6.69 | **$39** |
| Standard | $8.00 | **$70** |
| Express | $15.00 | — |

These are Shopify's starter defaults. Nobody set them for this catalog. Against the
real price ladder — $4, $5, $9, $14, $18, $22, $28, $38, $38, $45 — the $39 Economy
threshold is actively unprofitable.

### The arithmetic

Take the $38 Moon Hiker Wood Slice, which sits **one dollar** under the threshold.

| | Customer pays | Merchant nets after ~$6.69 postage |
|---|---|---|
| Buys the slice alone | $38 + $6.69 shipping = $44.69 | **$38.00** |
| Adds a $4 sticker → $42, free shipping | $42.00 | **$35.31**, minus sticker COGS |

The nudge that looks like an upsell — *"add $1 for free shipping"* — makes the
merchant **$2.69 worse off**, before the cost of the sticker. The threshold is set
below the point where absorbing shipping pays for itself.

It gets worse on heavy goods: the $45 Happy Stump Sign is 12in of pine and already
qualifies for free Economy. Real postage on it is very likely more than $6.69, and
the flat rate cannot tell the difference between a sign and a sticker.

### Fix

Two options, in order of preference:

1. **Weight-based rates.** Correct, and the reason `physical.weight_g` exists in the
   product schema. Currently null on all 8 products — so this is blocked until
   someone weighs them. That is a scale and ten minutes, not a project.
2. **Raise the Economy threshold** to ~$65–75 as an interim. A threshold should sit
   meaningfully above AOV, not one dollar above your most common mid-price item.

Do this **before** publishing the shipping policy, or you commit publicly to rates
you already intend to change.

### The one genuinely good nudge

Once thresholds are set correctly, the $38 price point is a gift: two products sit
exactly one dollar under a round threshold. A cart-drawer nudge at that boundary
converts unusually well. The mechanism is right — the threshold is just currently
in the wrong place.

## Fix order

Ranked by expected revenue per hour of work. Anything above a line is worthless
until the things below it are done.

| # | Fix | Why this rank | Blocked by |
|---|---|---|---|
| 0 | **Confirm payments are live** | `supportedDigitalWallets` is empty, which on a US store suggests no provider is activated. If true, *everything* below converts to $0. Not confirmable via API. | Settings → Payments, manual |
| 1 | **Create the three missing revenue lines** | Workshop seats, commission deposits, QR plaques. The business's own best products have no checkout path. | Nothing — this is buildable now |
| 2 | **Product photography** | For a craft brand whose entire proposition is *how the burned piece looks*, images are not decoration, they are the conversion mechanism. All 8 products have zero. | Camera, a morning |
| 3 | **Fix free-shipping thresholds** | Currently loses money per order. | Weigh the products, or set interim threshold |
| 4 | **Domain + store name** | `My Store 4` at a `myshopify.com` URL fails the trust check at exactly the moment a card comes out. | DNS + admin, manual |
| 5 | **Policies published** | Required at checkout; some payment providers demand them. | Paste from `../policies/` |
| — | *— everything above is plumbing; below is optimisation —* | | |
| 6 | CTA copy, cart nudges, POS bridge | Real levers, but they multiply a conversion rate that is currently structurally zero. | 0–5 |

## The POS ↔ online bridge

Worth designing now, cheap to build later. The booth is the strongest asset the
business has — Notion records the waitlist phone capture running at near-100%
because taking the number is part of the operation, and ~42 potential UGC pieces per
event.

The bridge that matters:

- **Booth QR → product page**, per product, carrying `?utm_source=booth`. The schema
  already has `surfaces.booth.qr_target` for exactly this.
- **Workshop waitlist → workshop product.** Today the number goes on a list and
  nothing automatic happens. Once a workshop product exists (fix #1), the list has
  somewhere to convert to.
- **UGC → product page.** Participants photograph their own pieces. Those photos are
  the product photography problem (#2) solving itself, if collected — `media.gallery`
  and `channels.social_posts[].ugc` in the schema are where they land.

## What is deliberately not recommended yet

- **CTA/copy A-B testing.** Requires traffic. The brand site has 0 email signups and
  1 workshop interest recorded; there is no traffic to split.
- **Discount codes.** Discounting an unphotographed product on an unbranded domain
  trains the wrong buyer and hides the real problem.
- **Abandoned-cart email.** Nothing can be abandoned until something can be bought.
