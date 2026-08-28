# Hosting the store on shop.doodle-wood.com

**Status: NOT DONE — blocked, requires manual action. This is the only step in the store build that an agent cannot perform.**

## Why this can't be automated

Two hard blockers, both verified rather than assumed:

1. **Shopify's Admin API has no mutation to attach a custom domain.** Domain connection
   is an Admin-UI-only operation. There is no `domainCreate`/`domainAdd` in the public
   GraphQL Admin API, so no token or scope makes this scriptable.
2. **DNS for `doodle-wood.com` is not on a connected provider.** The apex resolves to
   `185.158.133.1`, which is Hostinger address space. There is no Hostinger connector in
   this session, so the CNAME cannot be created programmatically either.

Current DNS state at time of writing:

```
doodle-wood.com        A       185.158.133.1     (live - brand site)
shop.doodle-wood.com   -       NXDOMAIN          (not configured)
```

## The two steps to complete it

Do these in order. Step 2 will fail verification until step 1 has propagated.

### Step 1 — Create the CNAME at Hostinger

In Hostinger → Domains → `doodle-wood.com` → DNS / Nameservers:

| Type  | Name   | Points to           | TTL  |
|-------|--------|---------------------|------|
| CNAME | `shop` | `shops.myshopify.com` | 3600 |

Do **not** point it at `000yqx-de.myshopify.com`. Shopify requires the shared
`shops.myshopify.com` target for custom domains; the shop-specific hostname will
not serve the storefront and will fail SSL provisioning.

Do **not** touch the apex `A` record — it serves the live brand site.

### Step 2 — Connect the domain in Shopify

Shopify Admin → Settings → Domains → **Connect existing domain** → enter
`shop.doodle-wood.com` → Verify connection.

Shopify then auto-provisions a Let's Encrypt certificate. This typically takes
a few minutes but is documented as up to 48 hours.

### Step 3 — Verify

```sh
getent hosts shop.doodle-wood.com          # should resolve
curl -sSI https://shop.doodle-wood.com | head -1   # expect HTTP/2 200
```

Then re-run `scripts/snapshot.sh` and commit, so the recorded state matches reality.

## Architecture note

Notion `40 · Storefront` describes a three-domain stack: `.com` (brand front door),
`.club` (maker network), `.shop` (commerce), and recommends building on subdomains
so a single login works across them. `shop.doodle-wood.com` as requested is
consistent with that recommendation — the `.shop` TLD, if acquired, should be
pointed at this subdomain as an alias rather than hosted separately.

## Do not skip

The brand site at `doodle-wood.com` is a **brand + story + email-capture** build with
no cart or checkout (per its own build prompt, recorded in Notion). Adding the Shopify
storefront on `shop.` is additive and does not disturb it — but any change to the apex
`A` record would take the live site down. Only add the `shop` CNAME.
