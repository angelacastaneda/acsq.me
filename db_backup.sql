PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL UNIQUE,
  file_name TEXT NOT NULL UNIQUE,
  content TEXT NOT NULL,
  description TEXT NOT NULL,
  pub_date TEXT NOT NULL CHECK(pub_date LIKE '____-__-__'),
  update_date TEXT NOT NULL CHECK(update_date LIKE '____-__-__'),
  thumbnail TEXT -- in json format but go engine can't handle real json
);
INSERT INTO posts VALUES(6,'a frosty stroll','a-frosty-stroll',replace('\n  <h2>a frosty stroll</h2>\n  <p><img src="https://preview.redd.it/hdnwdw6jy2f81.jpg?width=1024&amp;auto=webp&amp;v=enabled&amp;s=49a4b3cd118b59d5a59c92dd4a34740311ffa7ed" alt="photo of a woman walking beside a bright, cold street" /></p>\n','\n',char(10)),'I was walking over a bridge when I took this photo of this cool looking street. I didn’t notice the woman right away.','2022-01-31','2022-01-31','{"src":"https://preview.redd.it/hdnwdw6jy2f81.jpg?width=1024\u0026amp;auto=webp\u0026amp;v=enabled\u0026amp;s=49a4b3cd118b59d5a59c92dd4a34740311ffa7ed","alt":"photo of a woman walking beside a bright, cold street","title":""}');
INSERT INTO posts VALUES(8,'leaves that are green','leaves-that-are-green',replace('\n  <h2>{{ .Post.Title }}</h2>\n  <figure>\n    <img src="https://preview.redd.it/rgy06570s4w71.jpg?width=1024&auto=webp&v=enabled&s=6976a979db46bb48fa96e6fe61dcc6d1d2395141" alt="photo of tops of trees becoming orange" title="''I was twenty-one years when I wrote this song''">\n    <figcaption>Inspired by the <a href="https://youtu.be/OTlpdCPrx14" target="_blank" rel="noreferrer noopener">song</a></figcaption>\n  </figure>\n','\n',char(10)),'Took this pic on my way to university just by Lake Michigan.','2021-10-28','2021-10-28','{"src":"https://preview.redd.it/rgy06570s4w71.jpg?width=1024\u0026auto=webp\u0026v=enabled\u0026s=6976a979db46bb48fa96e6fe61dcc6d1d2395141","alt":"photo of tops of trees becoming orange","title":"''I was twenty-one years when I wrote this song''"}');
INSERT INTO posts VALUES(9,'what sank the titanic','what-sank-the-titanic',replace('\n  <h2>what sank the titanic</h2>\n  <figure>\n    <img src="https://preview.redd.it/jdoz0wwueid81.jpg?width=1024&amp;auto=webp&amp;v=enabled&amp;s=ddc97bed03ccca71f8559c82b3338ed8791ed9b4" alt="photo of a pretty chunk of ice" />\n  </figure>\n','\n',char(10)),'Took a pic of this interesting rocky/icy shoreline.','2022-01-23','2022-01-23','{"src":"https://preview.redd.it/jdoz0wwueid81.jpg?width=1024\u0026amp;auto=webp\u0026amp;v=enabled\u0026amp;s=ddc97bed03ccca71f8559c82b3338ed8791ed9b4","alt":"photo of a pretty chunk of ice","title":""}');
INSERT INTO posts VALUES(10,'natural log','natural-log',replace('\n  <h2>natural log</h2>\n  <figure>\n    <img src="https://preview.redd.it/zxw0z49sf2181.jpg?width=1024&amp;auto=webp&amp;v=enabled&amp;s=ec5260e1d2adf978777f16e535bdcc99fe592b84" title="ln(e^a) = a" alt="photo of bare, dead branch in front of goo pond" />\n  </figure>\n','\n',char(10)),'I went to a weird park with my roommate and took this pic of a branch there.','2021-11-21','2021-11-21','{"src":"https://preview.redd.it/zxw0z49sf2181.jpg?width=1024\u0026amp;auto=webp\u0026amp;v=enabled\u0026amp;s=ec5260e1d2adf978777f16e535bdcc99fe592b84","alt":"photo of bare, dead branch in front of goo pond","title":"ln(e^a) = a"}');
INSERT INTO posts VALUES(11,'new site... again','new-site-again',replace('\n  <h2>\n    new site... again\n  </h2>\n  \n  <p>\n    You think by this point I’d stop making such a hassle out of this, but no.\n    I completely scrapped my site I just made in the spring over the past two \n    months.\n  </p>\n  \n  <p>\n    The reason I destroyed it this time had nothing to do with bad design. Don’t get me wrong, it was kinda weird\n    looking for sure, but I enjoyed tinkering around with this tag or that. My main issue was I was expecting too much\n    from plain html files.\n  </p>\n\n  <h3>\n    the struggle\n  </h3>\n  \n  <p>\n    Most of my site design back in March and April was just vim (then eventually sed/regex) practice to see how\n    efficiently I could change the exact same header fifteen times if I ended up changing a tag or two. I wanted cool\n    features like article tagging and clever meta head stuff but by far the most cumbersome was multilingual content…\n  </p>\n  \n  <p>\n    This was definitely the biggest issue I had, so I wanna show my hack solution to dealing with it. Here’s how I’d\n    handle the html for my intro:\n  </p>\n  \n  <pre class="code"><code><strong>file: index.html</strong><hr>&lt;!doctype html&gt;\n&lt;html lang="en-US"&gt;\n  ...\n  &lt;h4 class="hl center"&gt;\n    &lt;span lang="en-US" xml:lang="en-US"&gt;Welcome to my site&lt;/span&gt;\n    &lt;span lang="es-US" xml:lang="es-US"&gt;Bienvenidos a mi sitio&lt;/span&gt;\n    &lt;span lang="de-DE" xml:lang="de-DE"&gt;Willkommen auf meiner Webseite&lt;/span&gt;\n  &lt;/h4&gt;\n  ...\n&lt;/html&gt;</code></pre>\n  \n  <p>\n    Then I had this rule in my css:\n  </p>\n  \n  <pre class="code"><code><strong>file: styles.css</strong><hr>/* what makes multilingual tags work */\nhtml:not([lang="en-US"]) :lang(en-US),\nhtml:not([lang="es-US"]) :lang(es-US),\nhtml:not([lang="de-DE"]) :lang(de-DE) {\n  display: none!important;\n}</code></pre>\n  \n  <p>\n    It was a pretty clever solution. I could just change the lang of the &lt;html&gt; tag itself as a switch for all my\n    spans. The source was a real mess, and chrome offering google translate would break the site, but it worked normally\n    besides those mishaps. The real issue was thinking this was a healthy way to write up a site.\n  </p>\n  \n  <p>\n    It didn’t take long before I burnt out writing it all by hand, so I stitched up html fragments into full files with a\n    shell script, but as long as I had to serve singular, static files, it was always gonna feel more like gluing a ransom\n    note than making a proper templating system.\n  </p>\n  \n  <figure>\n      <img src="/static/img/ransom_note.webp" alt="ransom note">\n      <figcaption>The culmination of my html scripting</figcaption>\n  </figure>\n\n  <h3>\n    a new way forward?\n  </h3>\n  \n  <p>\n    It was bad enough that I was just gonna move to some premade\n    <a href="https://gohugo.io/" target="_blank" rel="noopener noreferrer">hugo</a> solution when I found <a\n      href="https://j3s.sh/thought/my-website-is-one-binary.html" target="_blank" rel="noopener noreferrer">this\n      article</a>\n    that evangelized the power and simplicity of dynamic sites, specifically with Go.\n  </p>\n  \n  <p>\n    It seems as good a thing as any to learn, so I learned go’s funny syntactic slang like std, pkg, and fmt (no YHWH\n    yet.) It took very little time to figure out its power not just as a templating engine but also exposing to me what my\n    reverse proxy, <a href="https://nginx.org/" target="_blank" rel="noopener noreferrer">nginx</a>, was doing the whole\n    time serving particular domains by reading requests and writing responses.\n  </p>\n\n  <h3>\n    why dynamic is better\n  </h3>\n  \n  <p>\n    A good analogy that helped me understand the power of dynamic sites is cassette tapes. My weird systems I devised to\n    generate html files were me cutting up a bunch of riffs and drum solos and gluing them together into one long tape to\n    be read by a deck for some ’80s kid to enjoy.\n  </p>\n  \n  <p>\n    Dynamic sites are like those\n    <a href="https://www.youtube.com/shorts/_53MotXGiZc" target="_blank" rel="noopener noreferrer">digital to cassette\n      converters</a> people get for cars too old to have cd players or aux. They don’t really have any magnetic tape and\n    instead have a beam where the read head is to trick a cassette player into thinking it’s reading real tape, so you can\n    play anything arbitrarily.\n  </p>\n\n  <figure>\n    <img src="/static/img/cassette.jpg" alt="cassette tape held up in front of a car stereo system">\n    <figcaption>A static file, ready to serve</figcaption>\n  </figure>\n  \n  <p>\n    I do wanna show off the templating system (that hugo uses too) cause I think it’s real neat. Here’s a lang aware func\n    I made for the site in action:\n  </p>\n  \n  <pre class="code" id="demo"><code><strong>file: index.tmpl.html</strong><hr>&lt;!doctype html&gt;\n&lt;html lang="en-US"&gt;\n  ...\n  &lt;h4 class="hl center"&gt;\n    {{ "{{" }} translate .Lang\n    "Welcome to my site"\n    "Bienvenidos a mi sitio"\n    "Willkommen auf meiner Webseite" {{ "}}" }}\n  &lt;/h4&gt;\n  ...\n&lt;/html&gt;</code></pre>\n  \n  <p>\n    An here''s a live demo (+ dynamic lang buttons):\n  </p>\n\n  \n  <h4 class="hl center">\n    {{ translate .Lang\n    "Welcome to my site"\n    "Bienvenidos a mi sitio"\n    "Willkommen auf meiner Webseite"}}\n  </h4>\n\n  <nav class="center">\n    <strong>\n      <a href="{{.Scheme}}://{{translateHost "en-US" .Domain}}{{translatePath "en-US" .Path}}#demo">{{translateKeyword .Lang "english"}}</a>\n      <a href="{{.Scheme}}://{{translateHost "es-US" .Domain}}{{translatePath "es-US" .Path}}#demo">{{translateKeyword .Lang "spanish"}}</a>\n      <a href="{{.Scheme}}://{{translateHost "de-DE" .Domain}}{{translatePath "de-DE" .Path}}#demo">{{translateKeyword .Lang "german"}}</a>\n    </strong>\n  </nav>\n  \n  <h3>\n    simplicity is the best extensibility\n  </h3>\n\n  <p>\n    The coolest part isn’t even go''s http package or insane template engine. \n    Dynamic stuff such as having language specific urls is great, but it \n    wouldn''t have been possible if it didn''t make any sense to me.\n  </p>\n\n  <p>\n    The entire flow of control is self-written and, more importantly,\n    <em>self-understood</em>. Making your own handlers and router may seem hard\n    at first, but it''s so much easier than messing around with config files in\n    reverse proxies like nginx or apache you don''t understand in the first\n    place. Extensibility matters way less than comprehension. In fact, I''d argue \n    you can''t really utilize the former without the latter. Otherwise all you''re\n    gonna do is shrug your shoulders, copy a bunch of code, and pretend you\n    know what''s going on.\n  </p>\n\n  <h3>final thoughts</h3>\n\n  <p>\n    Static sites are cool. There''s a real charm in being able to set up a site\n    in three minutes and knowing that the url bar actually represents the\n    directory structure on the server itself. But they''re too restricive for me\n    to enjoy using them anymore. Writing my site in go was a fun way to stretch\n    out my feature set without having to use a godless flask/node/ruby on rails\n    webstack hosted on AWS/azure/google cloud.\n  </p>\n  \n  <p>\n    There’s still things I need to change like make a better md-&gt;html\n    system and move to a sqlite server cause constantly parsing files\n    for tags is painful, but after many months of writing angel-castaneda.com,\n    I finally feel like I’ve made something I can be proud of.\n  </p>\n  \n  <p>\n    Thanks again to\n    <a href="https://j3s.sh" target="_blank" rel="noopener noreferrer">jes</a> \n    for helping me to say no to big, clunky frameworks, dynamic or not. I also \n    wanna thank my friend <a href="https://www.alexscerba.com" target="_blank"\n    rel="noopener noreferrer">Alex</a> for helping me make my original site \n    back in 2021 as well as coding with me this summer to become \n    <em>dynamic</em>.\n  </p>\n','\n',char(10)),'Justifying why I spent two months remaking my site as a cover to learn go.','2023-08-07','2023-08-07','');
CREATE TABLE tags (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  category STRING NOT NULL DEFAULT 'content', -- for medium, content, and lang
  description TEXT
);
INSERT INTO tags VALUES(1,'articles','medium','Longform text written by me.');
INSERT INTO tags VALUES(2,'updates','content','For general updates to me or the site.');
INSERT INTO tags VALUES(3,'spanish','language','For my Spanish posts.');
INSERT INTO tags VALUES(4,'photos','content','Photos I''ve taken with a camera or phone or 3ds.');
INSERT INTO tags VALUES(5,'personal','content','For whenever I feel compelled to share something in my life online.');
INSERT INTO tags VALUES(6,'milwaukee','content','Where I get to talk about the coolest city in the midwest.');
INSERT INTO tags VALUES(7,'math','content','All posts tangentially related to math.');
INSERT INTO tags VALUES(8,'history','content','Anything cool related to the past I wanna talk about.');
INSERT INTO tags VALUES(9,'german','language','For my German posts.');
INSERT INTO tags VALUES(11,'english','language','{{ translate .Lang "For my English posts." "Mis entradas en inglés." "Meine Posten auf English." }}');
INSERT INTO tags VALUES(12,'code','content','{{ translate .Lang "For posts with code and computer stuff." "Para entradas sobre código." "Für Informatik." }}');
CREATE TABLE posts_tags (
  post_id INTEGER,
  tag_id INTEGER,
  PRIMARY KEY (post_id, tag_id),
  FOREIGN KEY (post_id) REFERENCES posts(id),
  FOREIGN KEY (tag_id) REFERENCES tags(id)
);
INSERT INTO posts_tags VALUES(5,4);
INSERT INTO posts_tags VALUES(5,6);
INSERT INTO posts_tags VALUES(5,11);
INSERT INTO posts_tags VALUES(6,4);
INSERT INTO posts_tags VALUES(6,6);
INSERT INTO posts_tags VALUES(6,11);
INSERT INTO posts_tags VALUES(7,4);
INSERT INTO posts_tags VALUES(7,6);
INSERT INTO posts_tags VALUES(7,11);
INSERT INTO posts_tags VALUES(8,4);
INSERT INTO posts_tags VALUES(8,6);
INSERT INTO posts_tags VALUES(8,11);
INSERT INTO posts_tags VALUES(9,4);
INSERT INTO posts_tags VALUES(9,6);
INSERT INTO posts_tags VALUES(9,11);
INSERT INTO posts_tags VALUES(10,4);
INSERT INTO posts_tags VALUES(10,11);
INSERT INTO posts_tags VALUES(10,7);
INSERT INTO posts_tags VALUES(11,12);
INSERT INTO posts_tags VALUES(11,1);
INSERT INTO posts_tags VALUES(11,2);
INSERT INTO posts_tags VALUES(11,11);
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('posts',11);
INSERT INTO sqlite_sequence VALUES('tags',12);
COMMIT;
