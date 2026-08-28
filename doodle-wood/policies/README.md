# Store policies

Ready-to-paste drafts for Shopify **Settings → Policies**.

## Why these are files and not already live

The Shopify MCP connection for this store lacks the `write_legal_policies` access
scope, so `shopPolicyUpdate` is denied:

```
Access denied for shopPolicyUpdate field.
Required access: `write_legal_policies` access scope.
```

Same story as the store name — Shopify's own docs say the shop resource
"doesn't let you update any information. Only the merchant can update this
information from inside the Shopify admin."

## How to apply

| File | Paste into |
|---|---|
| `refund-policy.html` | Settings → Policies → Refund policy |
| `shipping-policy.html` | Settings → Policies → Shipping policy |
| `privacy-policy.html` | Settings → Policies → Privacy policy |
| `terms-of-service.html` | Settings → Policies → Terms of service |
| `contact-information.html` | Settings → Policies → Contact information |

Each file has an HTML comment at the bottom listing what must be fixed before it
goes customer-facing. Those are not optional.

## Read this before publishing

These are **starter drafts written by an AI, not legal advice.** They are tuned to
the actual business — handmade variation, deposit-gated commissions, open-flame
workshops — which makes them a better starting point than Shopify's generic
templates, but they are a starting point.

Two sections carry real exposure and need a lawyer before you take money:

1. **Workshop liability** (`terms-of-service.html`). Paying attendees around hot
   tools and open flame. You very likely need a signed on-site waiver *as well as*
   this text, and your insurer may dictate the wording.
2. **Privacy** (`privacy-policy.html`). GDPR/CCPA exposure depends on where you
   sell. Compare against Shopify's auto-generated template and take the stronger
   wording of the two.

## Unresolved placeholders

These appear in the drafts and must be filled:

- **Contact email** — currently `junkmonkeystore@gmail.com`. Notion calls for a
  `doodle-wood.com` address; that mailbox does not exist yet.
- **Business address** — required by law in most jurisdictions and by some payment
  providers. Not invented here.
- **Legal entity name** — once registered.
- **Governing law / jurisdiction** — your state.
