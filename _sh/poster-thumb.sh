#!/bin/bash

# Small
SMALL_PNG="$1.small-thumb.png"
SMALL_WEBP="$1.small-thumb.webp"
rm -f $SMALL_WEBP
convert "$1" -resize 100x100 $SMALL_PNG
cwebp -q 90 $SMALL_PNG -o $SMALL_WEBP
rm $SMALL_PNG

# Medium
MEDIUM_PNG="$1.medium-thumb.png"
MEDIUM_WEBP="$1.medium-thumb.webp"
rm -f $MEDIUM_WEBP
convert "$1" -resize 700x700 $MEDIUM_PNG
cwebp -q 90 $MEDIUM_PNG -o $MEDIUM_WEBP
rm $MEDIUM_PNG