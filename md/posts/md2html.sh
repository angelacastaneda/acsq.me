#!/bin/bash
# converts md to .tmpl.html with pandoc

readonly POSTS_DIR="../../html/posts"

post=$(basename "$1" .md) 
echo "converting $post to HTML"

# pandoc for md to html
pandoc -f markdown-auto_identifiers -t html --template=post_template.tmpl.html --wrap=none --katex -o "$POSTS_DIR/$post.tmpl.html" "$post".md

# adds in filename to path
sed -i "s/posts\//posts\/$post/" "$POSTS_DIR/$post.tmpl.html"

# makes http & https links have target _blank and rel noopener noreferrer
perl -i -0pe 's/(<\W*a\W*[^>]*href=)(["'"'"']http[s]?:\/\/[^"'"'"'>]*["'"'"'])([^>]*>)/$1$2 target="_blank" rel="noopener noreferrer"$3/g' "$POSTS_DIR/$post.tmpl.html"

# fill in cool time thing at end
date=$(grep -oP '^date: \K\d{4}-\d{1,2}-\d{1,2}' "$post".md)
formatted_date=$(date -d "$date" +'%A, %B %d, %Y')
sed -i "s/><\/time/>$formatted_date<\/time/" "$POSTS_DIR/$post.tmpl.html"
