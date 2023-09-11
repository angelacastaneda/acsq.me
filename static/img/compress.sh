#!/bin/sh

res="1280x720"
qual="72"

og_img="$1"
pic_root="${og_img%.*}"

convert "$og_img" -resize "$res"\> comp_"$og_img"
convert comp_"$og_img" -quality "$qual"% comp_"$og_img"

cwebp comp_"$og_img" -o "$pic_root.webp"

rm -fv comp_"$og_img"

du -h "$pic_root"*
