---
title: new site... again
date: 2023-08-06
description: Justifying why I spent two months remaking my site as a cover to learn go.
tags: technology english articles
---

## new site... again

You think by this point I'd stop making such a hassle out of this, but no. I
completely scrapped my site I'd just made in the spring over the past two
months.

The reason I destroyed it this time had nothing to do with bad design. Don't
get me wrong, it was a kinda weird looking for sure, but I enjoyed tinkering
around with this tag or that. My main issue was I was expecting too much from
plain html files.

Most of my site design back in March and April was just vim (then eventually
sed/regex) practice to see how efficiently I could change the exact same header
fifteen times if I ended up changing a tag or two. I wanted cool features like
article tagging and clever meta head stuff but by far the most cumbersome was
multilingual content...

This was definitely the biggest issue I had, so I wanna show my hack
solution to dealing with it. Here's how I'd handle the html for my intro:

```html
<!doctype html>
<html lang="en-US">
⋮
  <h2 class="hl">
    <span lang="en-US" xml:lang="en-US">Welcome to my site</span>
    <span lang="es-US" xml:lang="es-US">Bienvenidos a mi sitio</span>
    <span lang="de-DE" xml:lang="de-DE">Willkommen auf meiner Webseite</span>
  </h2>
⋮
</html>
```

Then I'd have this rule in my css:

```css
/* what makes multilingual tags work */
html:not([lang="en-US"]) :lang(en-US),
html:not([lang="es-US"]) :lang(es-US),
html:not([lang="de-DE"]) :lang(de-DE) {
  display: none!important;
}
```

It was a pretty clever solution. I could just change the lang of the
&lt;html&gt; tag itself as a switch for all my spans. The source was a real
mess, and chrome offering google translate would break the site, but it worked
normally besides those mishaps. The real issue was thinking this was a healthy
way to write up a site.

It didn't take long before I burnt out writing it all by hand, so I stitched
up html fragments into full files with a shell script, but as long as I had to
serve singular, static files, it was always gonna feel more like gluing a
ransom note than making a proper templating system.

<style type="text/css">
#ransomizer {
  font-size:3em;
  line-height: normal;
  word-spacing:0.5em;
}

#ransomizer div {
  display:inline-block;
}

#ransomizer .rww {
  white-space: pre;
  display:inline;
  margin-left: .2em;
  margin-right: .2em;
}

#ransomizer .rr {
  -ms-transform: rotate(1.5deg);
  -webkit-transform: rotate(1.5deg);
  transform: rotate(1.5deg);
}

#ransomizer .rl {
  -ms-transform: rotate(-1.5deg);
  -webkit-transform: rotate(-1.5deg);
  transform: rotate(-1.5deg);
}

</style>
<a href="https://www.ransomizer.com">
  <div id="ransomizer"><div class="rww"><div class="rl" style="background-color:#0C8489;color:#fbffff;font-family:&#039;Verdana&#039;, Geneva, sans-serif; font-size:100%; font-weight:bold; background-image:url(https://i.imgur.com/1wxqouY.png) ; background-position: center center; box-shadow:1px -1px 2px #333; text-transform:lowercase; line-height:75%; margin:0.1em; padding:0.3em; vertical-align:0.1em; ">s</div><div class="rr" style="background-color:#F15770;color:#000000;font-family:&#039;Times New Roman&#039;, Times, serif; font-size:110%; background-image:url(https://i.imgur.com/1wxqouY.png) ; background-position: center top; box-shadow:1px -1px 2px #333; text-transform:lowercase; line-height:75%; margin:0.1em; padding:0.3em; vertical-align:-0.1em; ">t</div><div class="rl" style="background-color:#717732;color:#ffffff;font-family:&#039;Impact&#039;, Charcoal, sans-serif; font-size:100%; background-image:url(https://i.imgur.com/pwrAKPo.png) ; background-position: left center; box-shadow:1px 1px 2px #333; text-transform:lowercase; line-height:125%; margin:0.1em; padding:0.3em; vertical-align:-0.1em; ">a</div><div class="rl" style="background-color:#0C8489;color:#fbffff;font-family:&#039;Verdana&#039;, Geneva, sans-serif; font-size:100%; background-image:url(https://i.imgur.com/1wxqouY.png) ; background-position: left bottom; line-height:75%; margin:0.1em; padding:0.3em; vertical-align:0.1em; ">t</div><div class="rl" style="background-color:#006847;color:#ceffef;font-family:&#039;Courier&#039;, monospace; font-size:100%; font-weight:bold; background-image:url(https://i.imgur.com/pwrAKPo.png) ; background-position: center top; font-style:italic; text-transform:uppercase; line-height:75%; text-decoration:underline; margin:0.1em; padding:0.2em; vertical-align:0.1em; ">i</div><div class="rr" style="background-color:#D2A567;color:#000000;font-family:&#039;Verdana&#039;, Geneva, sans-serif; font-size:90%; font-weight:lighter; background-image:url(https://i.imgur.com/ruhP2kd.png) ; background-position: right bottom; font-style:italic; box-shadow:1px -1px 2px #333; text-transform:uppercase; line-height:75%; margin:0.1em; padding:0.2em; vertical-align:0.1em; ">c</div></div> <div class="rww"><div class="rr" style="background-color:#F8C83C;color:#000000;font-family:&#039;Trebuchet MS&#039;, Helvetica, sans-serif; font-size:110%; font-weight:lighter; background-image:url(https://i.imgur.com/ruhP2kd.png) ; background-position: right center; box-shadow:1px -1px 2px #333; text-transform:uppercase; line-height:100%; margin:0.1em; padding:0em; vertical-align:baseline; ">s</div><div class="rr" style="background-color:#006847;color:#ceffef;font-family:&#039;Courier&#039;, monospace; font-size:110%; font-weight:bold; background-image:url(https://i.imgur.com/pwrAKPo.png) ; background-position: center center; font-style:italic; box-shadow:1px 1px 2px #333; text-transform:uppercase; line-height:75%; margin:0.1em; padding:0.1em; vertical-align:0.1em; ">u</div><div class="rl" style="background-color:#803F1D;color:#ffffff;font-family:&#039;Impact&#039;, Charcoal, sans-serif; font-size:90%; background-image:url(https://i.imgur.com/ruhP2kd.png) ; background-position: center center; font-style:italic; box-shadow:-1px -1px 2px #333; line-height:75%; margin:0.1em; padding:0.3em; vertical-align:-0.1em; ">c</div><div class="rr" style="background-color:#0C8489;color:#fbffff;font-family:&#039;Courier&#039;, monospace; font-size:100%; font-weight:bolder; font-style:italic; box-shadow:-1px -1px 2px #333; text-transform:lowercase; line-height:125%; margin:0.1em; padding:0.2em; vertical-align:baseline; ">k</div><div class="rl" style="background-color:#006847;color:#ceffef;font-family:&#039;Times New Roman&#039;, Times, serif; font-size:90%; font-style:italic; text-transform:lowercase; line-height:100%; margin:0.1em; padding:0.3em; vertical-align:baseline; ">s</div></div></div>
</a>

It was bad enough that I was just gonna move to some premade
[hugo](https://gohugo.io/) solution when I found [this
article](https://j3s.sh/thought/my-website-is-one-binary.html) that evangelized
the power and simplicity of dynamic sites, specifically with Go.

It seems as good a thing as any to learn, so I learned go's funny syntactic
slang like std, pkg, and fmt (no YHWH yet.) It took very little time to
figure out its power not just as a templating engine but also exposing me
what my reverse proxy, nginx, was doing the whole time serving particular
domains by reading requests and writing responses.

A good analogy that helped me understand the power of dynamic sites is cassette
tapes. My weird systems I devised to generate html files were me cutting up a
bunch of riffs and drum solos and gluing them together into one long tape to be
read by a deck for some '80s kid to enjoy.

Dynamic sites are like those [digital to cassette
converters](https://www.youtube.com/shorts/_53MotXGiZc) people get for cars too
old to have cd players or aux. They don't actually have any magnetic tape and
instead have a beam where the read head is to trick a cassette player into
thinking it's reading real tape, so you can play anything arbitrarily.

I do wanna show off the templating system (that hugo uses too) cause I think
it's real neat. Here's a lang aware func I made for the site:

```html
<!doctype html>
<html lang="en-US">
⋮
  <h2 class="hl">
    \{\{ translate .Lang
    "Welcome to my site"
    "Bienvenidos a mi sitio"
    "Willkommen auf meiner Webseite"\}\}
  </h2>
⋮
</html>
```

An here it is in action: (lang buttons on the bottom)

<p class="center">
  {{ translate .Lang
  "Welcome to my site"
  "Bienvenidos a mi sitio"
  "Willkommen auf meiner Webseite"}}
</p>

The biggest thing for me isn't even the extensibility. Having language specific
urls is nice, but it was far more important to have something that made sense
to me without having to just shrug and pretend I know what it means.

There's still some things I'm gonna change like make a better md-&gt;html system
and probably move to a sqlite server cause constantly parsing files for tags is
painful, but after many months of writing angel-castaneda.com, I finally feel
like I've made something I can be really proud of.

Thanks again to [jes](https://j3s.sh) for helping me to say no to big, clunky
frameworks, dynamic or not. I also wanna thank my friend
[Alex](https://www.alexscerba.com) for helping me make my original site back in
2021 and suffering with me this summer to become: *dynamic*.
