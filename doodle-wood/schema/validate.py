#!/usr/bin/env python3
"""Validate Doodle Wood product records and report completeness per surface.

The schema is deliberately sparse-tolerant, so "valid" is a low bar. The useful
output here is the second half: which surfaces each product is actually ready to
publish to, and what is missing.

Usage:
    ./validate.py                    # validate everything under examples/ and ../catalog/products/
    ./validate.py path/to/rec.json   # validate specific records

Exit code 1 if any record fails schema validation. Missing fields are reported
but do NOT fail the run — incomplete is expected, invalid is not.
"""

import json
import pathlib
import sys

try:
    import jsonschema
except ImportError:
    sys.exit("Needs jsonschema:  pip install jsonschema")

HERE = pathlib.Path(__file__).parent
SCHEMA_PATH = HERE / "product.schema.json"

# What each surface needs before a product can go live on it.
SURFACE_REQUIREMENTS = {
    "shop": [
        "title",
        "handle",
        "commerce.price",
        "media.hero_image",
        "media.alt_text",
        "physical.weight_g",
        "narrative.description_short",
    ],
    "brand": [
        "title",
        "narrative.one_liner",
        "media.hero_image",
        "surfaces.brand.deep_link_to_shop",
    ],
    "booth": ["title", "commerce.booth_price", "physical.safety_notes"],
    "wholesale": ["title", "commerce.wholesale_price", "commerce.wholesale_moq"],
}


def get_path(obj, dotted):
    """Walk a dotted path; return None if any hop is missing."""
    cur = obj
    for part in dotted.split("."):
        if not isinstance(cur, dict):
            return None
        cur = cur.get(part)
        if cur is None:
            return None
    return cur


def main():
    schema = json.loads(SCHEMA_PATH.read_text())
    validator = jsonschema.Draft202012Validator(schema)

    if len(sys.argv) > 1:
        paths = [pathlib.Path(p) for p in sys.argv[1:]]
    else:
        paths = sorted(HERE.glob("examples/*.json"))
        paths += sorted((HERE.parent / "catalog" / "products").glob("*.json"))

    if not paths:
        print("No records found.")
        return 0

    failed = False
    for path in paths:
        record = json.loads(path.read_text())
        # Strip our own annotation keys before validating.
        record = {k: v for k, v in record.items() if not k.startswith("$") and k != "_note"}

        errors = sorted(validator.iter_errors(record), key=lambda e: e.path)
        name = record.get("id") or path.stem

        if errors:
            failed = True
            print(f"\nFAIL  {name}  ({path})")
            for err in errors:
                loc = ".".join(str(p) for p in err.path) or "(root)"
                print(f"        {loc}: {err.message}")
            continue

        print(f"\nOK    {name}")

        # Completeness per surface.
        for surface, required in SURFACE_REQUIREMENTS.items():
            missing = [f for f in required if get_path(record, f) in (None, "", [])]
            if missing:
                print(f"        {surface:<10} NOT READY  missing: {', '.join(missing)}")
            else:
                print(f"        {surface:<10} ready")

        # Margin check, only when both numbers exist.
        price = get_path(record, "commerce.price")
        cogs = get_path(record, "production.cogs")
        if price and cogs:
            margin = (price - cogs) / price * 100
            flag = "  <-- LOW" if margin < 50 else ""
            print(f"        margin     {margin:.0f}%{flag}")

        declared = set(record.get("_gaps", []))
        actual = {
            f
            for req in SURFACE_REQUIREMENTS.values()
            for f in req
            if get_path(record, f) in (None, "", [])
        }
        undeclared = actual - declared
        if undeclared:
            print(f"        undeclared gaps: {', '.join(sorted(undeclared))}")

    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main())
