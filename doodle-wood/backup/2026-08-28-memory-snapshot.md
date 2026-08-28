# Backup memory snapshot — 2026-08-28

Point-in-time record of what was true and what was decided, so a future session
(or a future you) can resume without re-deriving any of it.

## What this session changed in the live store

- Created 4 collections: `hand-burned-wood-slices`, `pendants-wearables`,
  `signs-wall-art`, `stickers-camp-goods`.
- Assigned all 8 existing products into them.
- Moved all 8 products from DRAFT → ACTIVE. They now have public product URLs.
- Fixed a bad auto-generated handle (`pendants-amp-wearables` → `pendants-wearables`)
  caused by HTML-escaping the ampersand in the collection title.

Nothing was deleted. Every change is reversible: set products back to DRAFT, delete
the 4 collections.

## What was found, not changed

- The store was **not** empty. A prior session (2026-08-27) had already created the
  8 products as drafts. This session structured and published them rather than
  recreating anything.
- Store is still named "My Store 4" — the Shopify default.
- Storefront password protection is **off**, so the store has been publicly reachable
  the whole time.
- All 8 products have **zero images**.

## Decisions taken, with reasoning

1. **Published products despite missing images.** Reversible in one call, and the store
   cannot reach customers under the brand domain anyway (see 2), so the exposure is a
   myshopify URL nobody has been given. Publishing was required for the store to be
   testable at all. Flagged as blocker #1 rather than silently accepted.
2. **Did not connect `shop.doodle-wood.com`.** Verified impossible from here: Shopify's
   Admin API exposes no domain-attach mutation, and `doodle-wood.com` DNS sits at
   Hostinger (apex → `185.158.133.1`) with no connector available. Documented as a
   two-step manual job in `docs/DOMAIN-SETUP.md` instead of faking completion.
3. **Built collections from real inventory, not from Notion's aspirational list.**
   Notion names "Festival Totems" and "Oracle Cards & Adventure Maps"; no such products
   exist. Creating empty collections would have looked like progress and been noise.
   The gap is recorded in `catalog/collections.json` under `_planned_but_not_built`.
4. **Chose `shop.` subdomain over the `.shop` TLD** that Notion `40 · Storefront` plans.
   The same Notion doc recommends building on subdomains so one login works across
   properties, so this follows its own advice and matches what was asked for.

## Where the code lives

`mon-key-mike/mcp-gateway`, branch `claude/doodle-wood-shopify-setup-t0n8di`,
directory `doodle-wood/`.

**This is a strange home for it.** That repo is a fork of Docker's Go MCP gateway and
has nothing to do with Doodle Wood. It was the only repository in scope for this
session and the branch was pre-assigned. A dedicated `doodle-wood` repo should be
created and this directory moved there — noted as an open item, not done unilaterally.

## Notion coordinates

- Company Library (wiki DB): `3c7b2c8e-c4cc-818d-aeb2-c8c4f763abe5`
- `40 · Storefront`: `3c7b2c8e-c4cc-8199-a849-c6d8dfa71261`
- Venture page `Doodle Wood`: `3c9b2c8e-c4cc-81de-90a0-ec90dad08b3d`
- Design System: `3c7b2c8e-c4cc-81aa-9bcf-d65bc606cb1e`
