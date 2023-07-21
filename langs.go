package main

import (
	"html/template"
	"strings"
	"time"
	"unicode"
)

func fetchLang(url string) string {
  if strings.HasPrefix(url, "es.") {
    return "es-US"
  }
  
  if strings.HasPrefix(url, "de.") {
    return "de-DE"
  }

  return "en-US"
}

func lastOne(index int, size int) bool{
  return index == size - 1
}

func translate(lang, en, es, de string) template.HTML {
  switch lang {
  case "es-US":
    return template.HTML(es)
  case "de-DE":
    return template.HTML(de)
  default:
    return template.HTML(en)
  }
}

func translateKeyword(lang, keyword string) string {
  keywordDictionary := map[string]map[string]string{
    /////////////// PAGES //////////////////// 
    "home": {
      "en-US": "home",
      "es-US": "inicio",
      "de-DE": "start",
    },
    "about": {
      "en-US": "about",
      "es-US": "conóceme",
      "de-DE": "über",
    },
    "posts": {
      "en-US": "posts",
      "es-US": "entradas",
      "de-DE": "posten",
    },
    "friends": {
      "en-US": "friends",
      "es-US": "amigos",
      "de-DE": "freunde",
    },
    "library": {
      "en-US": "library",
      "es-US": "biblioteca",
      "de-DE": "bibliothek",
    },
    "tags": {
      "en-US": "tags",
      "es-US": "etiquetas",
      "de-DE": "stichwörter", 
    },
    "contact": {
      "en-US": "contact",
      "es-US": "contacto",
      "de-DE": "kontakt", 
    },
    "todo": {
      "en-US": "todo",
      "es-US": "pendiente",
      "de-DE": "aufgaben", 
    },
    //////////////// TAGS ///////////////////////
    // medium
    "articles": {
      "en-US": "articles",
      "es-US": "artículos",
      "de-DE": "artikel",
    },
    "photos": {
      "en-US": "photos",
      "es-US": "fotos",
      "de-DE": "fotos",
    },
    // lang
    "english": {
      "en-US": "english",
      "es-US": "inglés",
      "de-DE": "englisch",
    },
    "spanish": {
      "en-US": "spanish",
      "es-US": "español",
      "de-DE": "spanisch",
    },
    "german": {
      "en-US": "german",
      "es-US": "alemán",
      "de-DE": "deutsch",
    },
    // tags
    "math": {
      "en-US": "math",
      "es-US": "matemáticas",
      "de-DE": "mathe",
    },
    "milwaukee": {
      "en-US": "milwaukee",
      "es-US": "milwaukee",
      "de-DE": "milwaukee",
    },
    "history": {
      "en-US": "history",
      "es-US": "historia",
      "de-DE": "geschichte",
    },
    "technology": {
      "en-US": "technology",
      "es-US": "tecnologia",
      "de-DE": "technologie",
    },
    "personal": {
      "en-US": "personal",
      "es-US": "personal",
      "de-DE": "persönliches",
    },
    ///////////////// TIME ////////////////////
    // days of week
    "monday": {
      "en-US": "monday",
      "es-US": "lunes",
      "de-DE": "montag",
    },
    "tuesday": {
      "en-US": "tuesday",
      "es-US": "martes",
      "de-DE": "dienstag",
    },
    "wednesday": {
      "en-US": "wednesday",
      "es-US": "miércoles",
      "de-DE": "mittwoch",
    },
    "thursday": {
      "en-US": "thursday",
      "es-US": "jueves",
      "de-DE": "donnerstag",
    },
    "friday": {
      "en-US": "friday",
      "es-US": "viernes",
      "de-DE": "freitag",
    },
    "saturday": {
      "en-US": "saturday",
      "es-US": "sábado",
      "de-DE": "samstag",
    },
    "sunday": {
      "en-US": "sunday",
      "es-US": "domingo",
      "de-DE": "sonntag",
    },
    // months
    "january": {
      "en-US": "january",
      "es-US": "enero",
      "de-DE": "januar",
    },
    "february": {
      "en-US": "february",
      "es-US": "febrero",
      "de-DE": "februar",
    },
    "march": {
      "en-US": "march",
      "es-US": "marzo",
      "de-DE": "märz",
    },
    "april": {
      "en-US": "april",
      "es-US": "abril",
      "de-DE": "april",
    },
    "may": {
      "en-US": "may",
      "es-US": "mayo",
      "de-DE": "mai",
    },
    "june": {
      "en-US": "june",
      "es-US": "junio",
      "de-DE": "juni",
    },
    "july": {
      "en-US": "july",
      "es-US": "julio",
      "de-DE": "juli",
    },
    "august": {
      "en-US": "august",
      "es-US": "agosto",
      "de-DE": "august",
    },
    "september": {
      "en-US": "september",
      "es-US": "septiembre",
      "de-DE": "september",
    },
    "october": {
      "en-US": "october",
      "es-US": "octubre",
      "de-DE": "oktober",
    },
    "november": {
      "en-US": "november",
      "es-US": "noviembre",
      "de-DE": "november",
    },
    "december": {
      "en-US": "december",
      "es-US": "diciembre",
      "de-DE": "dezember",
    },
  } 

  translation, ok := keywordDictionary[strings.ToLower(keyword)][lang] 
  if ok {
    if keyword == strings.ToUpper(keyword) {
      return strings.ToUpper(translation)
    } else if keyword == strings.Title(strings.ToLower(keyword)) {
      return strings.Title(translation)
    } else {
      return translation
    } 
  }

  return keyword
}

func anglicize(foreign string) string {
  ellisIsland := map[rune]string {
    'á': "a",
    'ä': "ae",
    'é': "e",
    'í': "i",
    'ñ': "n",
    'ó': "o",
    'ö': "oe",
    'ß': "ss",
    'ú': "u",
    'ü': "ue",
  }
  
  var domestic strings.Builder

  for _, char := range foreign {
    if char > unicode.MaxASCII {
      replaceChar := ellisIsland[char]  // if rune not in map, it just gets deleted
      domestic.WriteString(replaceChar)
    } else {
      domestic.WriteRune(char)
    }
  }
  return domestic.String()
}

func translateURL(lang, originalURL string) string {
  var translatedURL string
  urlPath := strings.Split(originalURL, "/")

  for _, urlSlice := range urlPath[1:] {
    switch urlSlice { // todo make something that lets keys be equal to each other to go directly from freunde -> amigos without friends middleman
    case "math","matematicas","mathe":
      urlSlice = anglicize(translateKeyword(lang, "math"))
    case "articles","articulos","artikel":
      urlSlice = anglicize(translateKeyword(lang, "articles"))
    case "photos","fotos":
      urlSlice = anglicize(translateKeyword(lang, "photos"))
    case "english","ingles","englisch":
      urlSlice = anglicize(translateKeyword(lang, "english"))
    case "spanish","espanol","spanisch":
      urlSlice = anglicize(translateKeyword(lang, "spanish"))
    case "german","aleman","deutsch":
      urlSlice = anglicize(translateKeyword(lang, "german"))
    case "milwaukee":
      urlSlice = anglicize(translateKeyword(lang, "milwaukee"))
    case "history","historia","geschichte":
      urlSlice = anglicize(translateKeyword(lang, "history"))
    case "technology","tecnologia","technologie":
      urlSlice = anglicize(translateKeyword(lang, "technology"))
    case "personal","personliches":
      urlSlice = anglicize(translateKeyword(lang, "personal"))
    case "about","conoceme","ueber":
      urlSlice = anglicize(translateKeyword(lang, "about"))
    case "posts","entradas","posten":
      urlSlice = anglicize(translateKeyword(lang, "posts"))
    case "friends","amigos","freunde":
      urlSlice = anglicize(translateKeyword(lang, "friends"))
    case "library","biblioteca","bibliothek":
      urlSlice = anglicize(translateKeyword(lang, "library"))
    case "tags","etiquetas","stichwoerter":
      urlSlice = anglicize(translateKeyword(lang, "tags"))
    case "todo","pendiente","aufgaben":
      urlSlice = anglicize(translateKeyword(lang, "todo"))
    default:
      urlSlice = anglicize(translateKeyword(lang, urlSlice))
    }

    translatedURL = translatedURL + "/" + urlSlice
  }

  return translatedURL
}

func translateDate(lang, iso string) (string, error) {
  date, err := time.Parse("2006-01-02", iso)
  if err != nil {
    return "", err
  }

  formattedDate := date.Format("Monday, January 2, 2006")
  dayOfWeek := date.Format("Monday")
  month := date.Format("January")
  switch lang {
  case "es-US":
    formattedDate = date.Format("Monday, 2 de January del 2006")
    diaDeLaSemana := translateKeyword("es-US", dayOfWeek)
    formattedDate = strings.ReplaceAll(formattedDate, dayOfWeek, diaDeLaSemana)
    mes := translateKeyword("es-US", strings.ToLower(month))
    formattedDate = strings.ReplaceAll(formattedDate, month, mes)
    return formattedDate, nil
  case "de-DE":
    formattedDate = date.Format("Monday, 02.01.2006")
    tagDerWoche := translateKeyword("de-DE", dayOfWeek)
    formattedDate = strings.ReplaceAll(formattedDate, dayOfWeek, tagDerWoche)
    return formattedDate, nil
  default:
    return formattedDate, nil
  }
}
