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

var dictionary = buildDictionary([]map[string]string{
  //////////////// PAGES ///////////////////////
  {
    "en-US": "home",
    "es-US": "inicio",
    "de-DE": "start",
  }, {
    "en-US": "about",
    "es-US": "conóceme", // non ASCII 
    "de-DE": "über", // non ASCII 
  }, {
    "en-US": "posts",
    "es-US": "entradas",
    "de-DE": "posten",
  }, {
    "en-US": "friends",
    "es-US": "amigos",
    "de-DE": "freunde",
  }, {
    "en-US": "library",
    "es-US": "biblioteca",
    "de-DE": "bibliothek",
  }, {
    "en-US": "tags",
    "es-US": "etiquetas",
    "de-DE": "stichwörter",  // non ASCII 
  }, {
    "en-US": "contact",
    "es-US": "contacto",
    "de-DE": "kontakt", 
  }, {
    "en-US": "todo",
    "es-US": "pendiente",
    "de-DE": "aufgaben", 
  //////////////// TAGS ///////////////////////
  }, { // medium
    "en-US": "articles",
    "es-US": "artículos", // non ASCII 
    "de-DE": "artikel",
  }, {
    "en-US": "photos",
    "es-US": "fotos",
    "de-DE": "fotos",
  }, { // lang
    "en-US": "english",
    "es-US": "inglés", // non ASCII 
    "de-DE": "englisch",
  }, {
    "en-US": "spanish",
    "es-US": "español", // non ASCII 
    "de-DE": "spanisch",
  }, {
    "en-US": "german",
    "es-US": "alemán", // non ASCII 
    "de-DE": "deutsch",
  }, { // tags
    "en-US": "math",
    "es-US": "matemáticas", // non ASCII 
    "de-DE": "mathe",
  }, {
    "en-US": "milwaukee",
    "es-US": "milwaukee",
    "de-DE": "milwaukee",
  }, {
    "en-US": "history",
    "es-US": "historia",
    "de-DE": "geschichte",
  }, {
    "en-US": "technology",
    "es-US": "tecnologia",
    "de-DE": "technologie",
  }, {
    "en-US": "personal",
    "es-US": "personal",
    "de-DE": "persönliches", // non ASCII 
  ///////////////// TIME ////////////////////
  }, { // days of week
    "en-US": "monday",
    "es-US": "lunes",
    "de-DE": "montag",
  }, {
    "en-US": "tuesday",
    "es-US": "martes",
    "de-DE": "dienstag",
  }, {
    "en-US": "wednesday",
    "es-US": "miércoles", // non ASCII 
    "de-DE": "mittwoch",
  }, {
    "en-US": "thursday",
    "es-US": "jueves",
    "de-DE": "donnerstag",
  }, {
    "en-US": "friday",
    "es-US": "viernes",
    "de-DE": "freitag",
  }, {
    "en-US": "saturday",
    "es-US": "sábado", // non ASCII 
    "de-DE": "samstag",
  }, {
    "en-US": "sunday",
    "es-US": "domingo",
    "de-DE": "sonntag",
  }, { // months
    "en-US": "january",
    "es-US": "enero",
    "de-DE": "januar",
  }, {
    "en-US": "february",
    "es-US": "febrero",
    "de-DE": "februar",
  }, {
    "en-US": "march",
    "es-US": "marzo",
    "de-DE": "märz", // non ASCII 
  }, {
    "en-US": "april",
    "es-US": "abril",
    "de-DE": "april",
  }, {
    "en-US": "may",
    "es-US": "mayo",
    "de-DE": "mai",
  }, {
    "en-US": "june",
    "es-US": "junio",
    "de-DE": "juni",
  }, {
    "en-US": "july",
    "es-US": "julio",
    "de-DE": "juli",
  }, {
    "en-US": "august",
    "es-US": "agosto",
    "de-DE": "august",
  }, {
    "en-US": "september",
    "es-US": "septiembre",
    "de-DE": "september",
  }, {
    "en-US": "october",
    "es-US": "octubre",
    "de-DE": "oktober",
  }, {
    "en-US": "november",
    "es-US": "noviembre",
    "de-DE": "november",
  }, {
    "en-US": "december",
    "es-US": "diciembre",
    "de-DE": "dezember",
  },
})

func buildDictionary(maps []map[string]string) map[string]map[string]string {
  rosetta := make(map[string]map[string]string)
  for _, m := range maps {
    for lang, translation := range m {
      if rosetta[lang] == nil {
        rosetta[lang] = make(map[string]string)
      }
      for _, word2 := range m {
        rosetta[lang][word2] = translation
      }
    }
  } 
  return rosetta
}

func unAnglicize(word string) string {
  authenticDictionary := map[string]string{
    // español
    "articulos": "artículos",
    "conoceme": "conóceme",
    "matematicas": "matemáticas",
    "ingles": "inglés",
    "espanol": "español",
    "aleman": "alemán",
    "miercoles": "miércoles",
    "sabado": "sábado",
    // deutsch
    "ueber": "über",
    "persoenliches": "persönliches",
    "stichwoerter": "stichwörter",
    "maerz": "märz",
  }
  
  palabra, okey := authenticDictionary[word]
  if okey {
    return palabra
  }
  return word
}

func translateKeyword(lang, keyword string) string {
  translation, ok := dictionary[lang][unAnglicize(strings.ToLower(keyword))]
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

func translatePath(lang, originalURL string) string {
  var translatedURL string
  urlPath := strings.Split(originalURL, "/")

  for _, urlSlice := range urlPath[1:] {
    urlSlice = anglicize(translateKeyword(lang, urlSlice))
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

func translateHost(lang, domain string) string {
  subLangs := []string{"www.", "en.", "es.", "de."}
  // localhost:4000 -> de.localhost:4000
  // www.angel-castaneda.com -> es.angel-castaneda.com
  // angel.localhost:8080 -> en.angel-castaneda.com
  // en.angel.localhost:8080 -> www.angel.localhost:8080

  switch lang{
  case "en-US":
    for _, lang := range subLangs {
      if strings.HasPrefix(domain, lang) {
        return strings.ReplaceAll(domain, lang, "en.")
      }
    }
    return "en." + domain
  case "es-US":
    for _, lang := range subLangs {
      if strings.HasPrefix(domain, lang) {
        return strings.ReplaceAll(domain, lang, "es.")
      }
    }
    return "es." + domain
  case "de-DE":
    for _, lang := range subLangs {
      if strings.HasPrefix(domain, lang) {
        return strings.ReplaceAll(domain, lang, "de.")
      }
    }
    return "es." + domain
  default:
    for _, lang := range subLangs {
      if strings.HasPrefix(domain, lang) {
        return strings.ReplaceAll(domain, lang, "www.")
      }
    }
    return "www." + domain
  }
}
