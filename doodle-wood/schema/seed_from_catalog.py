#!/usr/bin/env python3
"""Seed schema-conformant product records from the verified catalog snapshot.

Creates one file per product under ../catalog/products/, filling only what is
actually known from Shopify and leaving everything else null with an honest
_gaps list. Does NOT invent narrative, cost, or weight data.

Safe to re-run: skips records that already exist so hand-written detail is never
clobbered. Pass --force to regenerate anyway.

    ./seed_from_catalog.py [--force]
"""

import json
import pathlib
import sys

HERE = pathlib.Path(__file__).parent
CATALOG = HERE.parent / "catalog" / "products.json"
OUT_DIR = HERE.parent / "catalog" / "products"

# Character name per product handle, from the product descriptions.
CHARACTERS = {
    "campfire-monkey-wood-slice": "Campfire Monkey",
    "moon-hiker-wood-slice": "Moon Hiker",
    "banana-bandit-pendant": "Banana Bandit",
    "trail-monkey-patch": "Trail Monkey",
    "happy-stump-welcome-sign": "Happy Stump",
    "mushroom-buddy-sticker": "Mushroom Buddy",
    "cosmic-acorn-sticker": "Cosmic Acorn",
    "doodle-wood-camp-mug": None,  # logo product, no character
}

MOTIFS = {
    "monkey": ["monkey"],
    "mushroom": ["mushroom"],
    "acorn": ["flora"],
    "adventure": ["adventure"],
    "camp": ["camp"],
    "home": [],
    "sticker": [],
    "wearable": [],
    "sign": [],
    "pyrography": [],
    "mug": ["camp"],
    "doodle-wood": [],
}

COLLECTION_GIDS = {
    "hand-burned-wood-slices": "gid://shopify/Collection/337792368803",
    "pendants-wearables": "gid://shopify/Collection/337792401571",
    "signs-wall-art": "gid://shopify/Collection/337792434339",
    "stickers-camp-goods": "gid://shopify/Collection/337792467107",
}

# Fields every seeded record is missing until a human fills them in.
BASE_GAPS = [
    "media.hero_image",
    "media.alt_text",
    "media.source_doodle",
    "physical.weight_g",
    "production.cogs",
    "production.burn_minutes",
    "commerce.booth_price",
    "commerce.wholesale_price",
    "commerce.wholesale_moq",
    "narrative.origin_story",
    "narrative.description_long",
    "narrative.description_short",
    "physical.safety_notes",
    "surfaces.brand.deep_link_to_shop",
    "sources.notion_page",
]


def build(product):
    handle = product["handle"]
    collection = product["collection"]

    motif = []
    for tag in product["tags"]:
        motif.extend(MOTIFS.get(tag, []))
    motif = sorted(set(motif))

    return {
        "id": f"dw-{handle}",
        "sku": product["variants"][0]["sku"].rsplit("-", 1)[0]
        if len(product["variants"]) > 1
        else product["variants"][0]["sku"],
        "title": product["title"],
        "handle": handle,
        "status": "active" if product["status"] == "ACTIVE" else "draft",
        "version": "v1.0",
        "classification": {
            "product_type": product["product_type"],
            "collection": collection,
            "secondary_collections": [],
            "tags": product["tags"],
            "character": CHARACTERS.get(handle),
            "motif": motif,
            "gift_occasion": [],
        },
        "narrative": {
            "one_liner": None,
            "description_short": None,
            "description_long": None,
            "origin_story": None,
            "grain_note": None,
            "voice_warnings": [
                "never describe as food-safe",
                "avoid 'rustic' - Notion brand doc explicitly rejects it",
            ],
        },
        "physical": {
            "wood_species": None,
            "dimensions_mm": None,
            "weight_g": None,
            "finish": None,
            "burn_depth": None,
            "hardware": None,
            "one_of_one": None,
            "care_instructions": None,
            "safety_notes": None,
        },
        "production": {
            "make_mode": None,
            "burn_minutes": None,
            "blank_source": None,
            "blank_cost": None,
            "consumables_cost": None,
            "cogs": None,
            "batch_size": None,
            "tooling": [],
            "difficulty": None,
            "workshop_suitable": None,
        },
        "commerce": {
            "currency": "USD",
            "price": float(product["variants"][0]["price"]),
            "compare_at_price": None,
            "booth_price": None,
            "wholesale_price": None,
            "wholesale_moq": None,
            "margin_pct": None,
            "tax_code": None,
            "inventory_tracked": True,
            "inventory_on_hand": sum(v["inventory"] for v in product["variants"]),
        },
        "variants": [
            {
                "sku": v["sku"],
                "title": v["title"],
                "option_size": v["title"].split(" / ")[0] if " / " in v["title"] else None,
                "option_material": v["title"].split(" / ")[1] if " / " in v["title"] else None,
                "price": float(v["price"]),
                "weight_g": None,
                "inventory_on_hand": v["inventory"],
                "shopify_variant_gid": None,
            }
            for v in product["variants"]
        ],
        "media": {
            "hero_image": None,
            "gallery": [],
            "alt_text": None,
            "burn_video": None,
            "source_doodle": None,
            "photography_status": "none",
        },
        "surfaces": {
            "shop": {
                "listed": product["status"] == "ACTIVE",
                "url": f"https://000yqx-de.myshopify.com/products/{handle}",
                "shopify_product_gid": product["shopify_gid"],
                "shopify_collection_gid": COLLECTION_GIDS.get(collection),
            },
            "brand": {
                "featured": False,
                "url": None,
                "lovable_component": None,
                "deep_link_to_shop": None,
            },
            "booth": {"carried": None, "display_tier": None, "qr_target": None},
            "wholesale": {"offered": None, "line_sheet_ref": None},
            "etsy": {"listed": False, "listing_id": None, "url": None},
        },
        "channels": {
            "utm_default": "?utm_source=brand&utm_medium=referral&utm_campaign=doodle-wood",
            "short_link": None,
            "affiliate_links": [],
            "social_posts": [],
        },
        "seo": {
            "meta_title": None,
            "meta_description": None,
            "keywords": [],
            "canonical_url": None,
        },
        "sources": {
            "notion_page": None,
            "concept_register_ref": CHARACTERS.get(handle),
            "design_files": [],
            "cost_model_ref": None,
            "internal_docs": [],
        },
        "lifecycle": {
            "created_at": "2026-08-27",
            "updated_at": "2026-08-28",
            "last_verified_against_shopify": "2026-08-28",
            "confidence": "confirmed",
            "retired_reason": None,
        },
        "_gaps": sorted(BASE_GAPS + ["narrative.one_liner", "seo.meta_title"]),
    }


def main():
    force = "--force" in sys.argv
    catalog = json.loads(CATALOG.read_text())
    OUT_DIR.mkdir(parents=True, exist_ok=True)

    written = skipped = 0
    for product in catalog["products"]:
        record = build(product)
        path = OUT_DIR / f"{record['id']}.json"
        if path.exists() and not force:
            print(f"skip   {record['id']} (exists)")
            skipped += 1
            continue
        path.write_text(json.dumps(record, indent=2) + "\n")
        print(f"write  {record['id']}")
        written += 1

    print(f"\n{written} written, {skipped} skipped -> {OUT_DIR}")
    print("Now run ./validate.py to see per-surface readiness.")


if __name__ == "__main__":
    main()
