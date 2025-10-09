#!/bin/bash

# Script to remove and recreate Git tags

set -e  # Exit on error

# Define tags
TAGS=(
    "addons/base/v0.3.0"
    "addons/web/v0.3.0"
    "hexya/v0.3.0"
    "pool/v0.2.0"
)

echo "=== Removing local tags ==="
for tag in "${TAGS[@]}"; do
    if git tag -l "$tag" | grep -q "$tag"; then
        git tag -d "$tag"
        echo "✓ Deleted local tag: $tag"
    else
        echo "- Local tag not found: $tag"
    fi
done

echo ""
echo "=== Removing remote tags ==="
for tag in "${TAGS[@]}"; do
    if git ls-remote --tags origin | grep -q "refs/tags/$tag"; then
        git push origin ":refs/tags/$tag"
        echo "✓ Deleted remote tag: $tag"
    else
        echo "- Remote tag not found: $tag"
    fi
done

echo ""
echo "=== Creating new tags on current HEAD ==="
for tag in "${TAGS[@]}"; do
    git tag -a "$tag" -m "Release $tag"
    echo "✓ Created tag: $tag"
done

echo ""
echo "=== Pushing tags to remote ==="
git push origin --tags

echo ""
echo "=== Done! ==="
echo "All tags have been recreated on current HEAD and pushed to remote."