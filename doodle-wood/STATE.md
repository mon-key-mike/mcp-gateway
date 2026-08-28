# Doodle Wood store — state of record

**Last verified against the live store: 2026-08-28**
Everything below was read from the Shopify Admin API, not assumed.

## Store identity

| Field | Value | OK? |
|---|---|---|
| Shop name | `My Store 4` | **NO — Admin-UI-only, not settable via API** |
| myshopify domain | `000yqx-de.myshopify.com` | working |
| Target domain | `shop.doodle-wood.com` | **NOT CONNECTED** |
| Plan | Basic | ok |
| Currency | USD | ok |
| Contact email | `junkmonkeystore@gmail.com` | not a doodle-wood.com address |
| Password protection | **disabled — storefront is publicly reachable** | see risk below |
| Live theme | Horizon | ok (Crave, Tinker unpublished) |

## Catalog — built and live

8 products, all moved DRAFT → ACTIVE on 2026-08-28. 4 collections created.

| Collection | Products | Price range |
|---|---|---|
| Hand-Burned Wood Slices | 2 | $28–38 |
| Pendants & Wearables | 2 | $9–18 |
| Signs & Wall Art | 1 | $45 |
| Stickers & Camp Goods | 3 | $4–22 |

Full detail: `catalog/products.json`, `catalog/collections.json`.

## Blockers before this store can launch

Ranked. The first two are launch-blocking.

1. **Every product has zero images** (`mediaCount: 0` on all 8). A pyrography brand
   whose entire value proposition is *how the burned piece looks* cannot sell without
   photography. This is the single biggest gap.
2. **`shop.doodle-wood.com` is not connected.** Requires a DNS change at Hostinger and
   a manual step in Shopify Admin — neither is automatable. See `docs/DOMAIN-SETUP.md`.
3. **Store is named "My Store 4"** and shows that name in the browser tab, checkout,
   and every transactional email. Change in Settings → Store details.
4. **Storefront is public with no password.** Combined with (1) and (3), anyone who
   finds `000yqx-de.myshopify.com` right now sees an unbranded store with imageless
   products. Either add a password until launch, or fix 1–3 quickly.
5. **Policies are written but not published.** Drafts live in `policies/` and must be
   pasted into Settings → Policies by hand — the Shopify connection lacks the
   `write_legal_policies` scope. They contain deliberate placeholders (contact email,
   business address) that are customer-facing and must be filled first. The workshop
   liability and privacy sections need a lawyer before taking money.
6. **Shipping rates ARE configured** (corrected 2026-08-28 — an earlier note here
   wrongly said they were not). Flat rates on a Domestic zone: Economy $6.69 (free
   ≥ $39), Standard $8.00 (free ≥ $70), Express $15.00. These are Shopify defaults
   and the free-shipping thresholds are mispriced against this catalog — see
   `docs/FUNNEL.md`.
7. **Payment provider not confirmed.** `supportedDigitalWallets` is empty, which on a
   US store strongly suggests no provider is activated. Not confirmable via API —
   check Settings → Payments. Until this is live, every other funnel fix converts
   to $0.
8. **The three highest-value revenue lines have no product in the store.** Notion
   records $30 workshop seats as the highest-throughput item, deposit-gated
   commissions, and B2B QR plaques from $72 as the highest-value repeat item.
   None exists as a purchasable product. See `docs/FUNNEL.md`.
7. Notion `40 · Storefront` plans commerce on `doodle-wood.shop`; the request here is
   `shop.doodle-wood.com`. Reconciled in favour of the subdomain — consistent with the
   same doc's own recommendation to build on subdomains. Noted so it isn't re-litigated.

## Relationship to doodle-wood.com

`doodle-wood.com` is live and is a **brand + story + email-capture site with no cart or
checkout**, by its own build spec. This Shopify store is the commerce surface that site
never had. Adding `shop.` is additive — do not touch the apex `A` record
(`185.158.133.1`), which serves the live brand site.

Open item carried from Notion: `auth.users` on the brand site is empty, so no admin
exists and the captured email list is unreadable. Unrelated to this store, still blocking
that funnel.
