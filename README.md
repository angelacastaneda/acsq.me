# [acsq.me](https://www.acsq.me): personal website of ángel castañeda

I wrote my site in go to escape the comforts of nginx and actually learn how a
webserver works.

## features

* templating html

* real multilingual content/url localization for
  [spanish](https://es.acsq.me) and
  [german](https://de.acsq.me)

* [atom feed](https://www.acsq.me/atom.xml)

* [fancy error pages](https://www.acsq.me/page/doesnt/exist)

* md -> html for articles

* blog organized w/ tagging system

* gzipping

* sqlite server

## License

Code for this site is under the GPLv3. See [`LICENSE`](./LICENSE) for more
details.

This site uses my [dblog repo](https://git.acsq.me/dblog) as an interface
between the sqlite database and my code. That package is licensed under the
LGPLv3. Check it's [`LICENSE`](https://git.acsq.me/dblog/tree/LICENSE/)
