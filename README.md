# [acsq.me](https://www.acsq.me): personal website of ángel castañeda

I wrote my personal site in golang to escape the comforts of static nginx and
actually learn how a webserver works.

Original inspiration came from [this blog post](https://j3s.sh/thought/my-website-is-one-binary.html).

## features

* templating html files that are basically little jigsaw puzzles

* real multilingual content/url localization for
  [spanish](https://es.acsq.me) and
  [german](https://de.acsq.me)

* [atom feed](https://www.acsq.me/atom.xml) for rss readers

* fancy http [status error pages](https://www.acsq.me/page/doesnt/exist)

* md -> html for articles

* blog organized w/ tagging system from my [dblog
  repo](https://git.acsq.me/dblog)

* gzipping

## license

Code for this site is under the GPLv3. See [`LICENSE`](./LICENSE) for more
details.

This site uses my [dblog repo](https://git.acsq.me/dblog) as an interface
between the sqlite database and my code. That package is licensed under the
LGPLv3. See it's [`LICENSE`](https://git.acsq.me/dblog/tree/LICENSE/) file
here.
