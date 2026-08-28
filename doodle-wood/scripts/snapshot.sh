#!/usr/bin/env bash
# Pull live Shopify state into catalog/ so git tracks every change to the store.
#
# Usage:
#   export SHOPIFY_STORE=000yqx-de.myshopify.com
#   export SHOPIFY_ADMIN_TOKEN=shpat_...      # never commit this
#   ./scripts/snapshot.sh
#
# Then: git diff catalog/  -- shows exactly what changed in the store since last commit.

set -euo pipefail

API_VERSION="2026-07"
STORE="${SHOPIFY_STORE:?set SHOPIFY_STORE, e.g. 000yqx-de.myshopify.com}"
TOKEN="${SHOPIFY_ADMIN_TOKEN:?set SHOPIFY_ADMIN_TOKEN (Admin API access token)}"

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="$HERE/catalog"
STAMP="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

mkdir -p "$OUT"

gql() {
  curl -sS -X POST \
    "https://${STORE}/admin/api/${API_VERSION}/graphql.json" \
    -H "X-Shopify-Access-Token: ${TOKEN}" \
    -H "Content-Type: application/json" \
    --data @<(jq -n --arg q "$1" '{query: $q}')
}

echo "==> Snapshotting products from ${STORE}"
gql 'query {
  products(first: 250) {
    edges { node {
      id handle title productType status vendor tags mediaCount { count }
      variants(first: 100) { edges { node { sku title price inventoryQuantity } } }
    } }
  }
}' | jq --arg s "$STAMP" '{_snapshot_taken: $s, data: .data.products.edges | map(.node)}' \
   > "$OUT/live-products.json"

echo "==> Snapshotting collections"
gql 'query {
  collections(first: 250) {
    edges { node {
      id handle title sortOrder productsCount { count }
      products(first: 100) { edges { node { handle } } }
    } }
  }
}' | jq --arg s "$STAMP" '{_snapshot_taken: $s, data: .data.collections.edges | map(.node)}' \
   > "$OUT/live-collections.json"

echo "==> Snapshotting shop + domain state"
gql 'query {
  shop { name myshopifyDomain currencyCode primaryDomain { host sslEnabled } }
  onlineStore { passwordProtection { enabled } }
}' | jq --arg s "$STAMP" '{_snapshot_taken: $s, data: .data}' \
   > "$OUT/live-shop.json"

echo
echo "Snapshot written to catalog/live-*.json at ${STAMP}"
echo "Review drift, then commit:"
echo "    git add doodle-wood/catalog && git commit -m 'chore(doodle-wood): store snapshot ${STAMP}'"

# Flag the things that block a real launch.
MISSING_MEDIA=$(jq '[.data[] | select(.mediaCount.count == 0)] | length' "$OUT/live-products.json")
if [ "$MISSING_MEDIA" -gt 0 ]; then
  echo
  echo "WARNING: ${MISSING_MEDIA} product(s) still have no images. The storefront is not launch-ready."
  jq -r '.data[] | select(.mediaCount.count == 0) | "  - " + .title' "$OUT/live-products.json"
fi
