#!/bin/bash
# converts md to .tmpl.html with pandoc

post=$(basename $1 .md) 
echo "converting $post to HTML"

# pandoc for md to html
pandoc -f markdown-auto_identifiers -t html --template=post_template.tmpl.html --wrap=none --katex -o ../../html/posts/$post.tmpl.html $post.md

# adds in filename to path
sed -i "s/posts\//posts\/$post/" ../../html/posts/$post.tmpl.html

# makes http & https links have target _blank and rel noopener noreferrer
perl -i -0pe 's/(<\W*a\W*[^>]*href=)(["'"'"']http[s]?:\/\/[^"'"'"'>]*["'"'"'])([^>]*>)/$1$2 target="_blank" rel="noopener noreferrer"$3/g' ../../html/posts/$post.tmpl.html

# fill in cool time tag at end
date=$(grep -oP '^date: \K\d{4}-\d{1,2}-\d{1,2}' $post.md)
formatted_date=$(date -d $date +'%A, %B %d, %Y')
sed -i "s/><\/time/>$formatted_date<\/time/" ../../html/posts/$post.tmpl.html
