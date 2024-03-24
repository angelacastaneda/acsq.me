package main

import (
	"html/template"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func lastOne(index int, size int) bool {
	return index == size-1
}

func translate(lang, en, es, de string) string {
	switch lang {
	case "es-US":
		return es
	case "de-DE":
		return de
	default:
		return en
	}
}

func translateHTML(lang, en, es, de string) template.HTML {
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
		"de-DE": "über",     // non ASCII
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
		"en-US": "recommend",
		"es-US": "recomendar",
		"de-DE": "empfehlen",
	}, {
		"en-US": "tags",
		"es-US": "etiquetas",
		"de-DE": "stichwörter", // non ASCII
	}, {
		"en-US": "contact",
		"es-US": "contacto",
		"de-DE": "kontakt",
	}, {
		"en-US": "todo",
		"es-US": "pendiente",
		"de-DE": "aufgaben",
		//////////////// MISC ///////////////////////
	}, {
		"en-US": "hi",
		"es-US": "hola",
		"de-DE": "hallo",
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
		"en-US": "code",
		"es-US": "código", // non ASCII
		"de-DE": "code",
	}, {
		"en-US": "updates",
		"es-US": "actualizaciones",
		"de-DE": "aktualisierungen",
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

// this awesome func come from a stack overflow post, and is not my own
// https://stackoverflow.com/a/76735864/21316874 ; it falls under CC BY-SA 4.0
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
		"articulos":   "artículos",
		"codigo":      "código",
		"conoceme":    "conóceme",
		"matematicas": "matemáticas",
		"ingles":      "inglés",
		"espanol":     "español",
		"aleman":      "alemán",
		"miercoles":   "miércoles",
		"sabado":      "sábado",
		// deutsch
		"ueber":         "über",
		"persoenliches": "persönliches",
		"stichwoerter":  "stichwörter",
		"maerz":         "märz",
	}

	palabra, okey := authenticDictionary[word]
	if okey {
		return palabra
	}
	return word
}

func translateKeyword(lang, keyword string) string {
	fmtedKeyword := unAnglicize(strings.ToLower(keyword))
	translation, ok := dictionary[lang][fmtedKeyword]

	if !ok {
		return keyword
	}

	caser := cases.Title(language.Und) // TODO learn to compose spanish, german, and english
	if keyword == strings.ToUpper(keyword) {
		return strings.ToUpper(translation)
	} else if keyword == caser.String(strings.ToLower(keyword)) {
		return caser.String(translation)
	}
	return translation
}

func anglicize(foreign string) string {
	ellisIsland := map[rune]string{
		// deutsch
		'ä': "ae",
		'ö': "oe",
		'ß': "ss",

		// espanol
		'á': "a",
		'é': "e",
		'í': "i",
		'ñ': "n",
		'ó': "o",
		'ú': "u",

		// both
		'ü': "ue",
	}

	var domestic strings.Builder

	for _, char := range foreign {
		if char > unicode.MaxASCII {
			replaceChar := ellisIsland[char] // if rune not in map, it just gets deleted
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

	last := urlPath[len(urlPath)-1]
	addHTML := false
	if strings.HasSuffix(last, ".html") {
		addHTML = true
		urlPath[len(urlPath)-1] = strings.TrimSuffix(last, ".html")
	}

	for _, urlSlice := range urlPath[1:] {
		urlSlice = anglicize(translateKeyword(lang, urlSlice))
		translatedURL = translatedURL + "/" + urlSlice
	}

	if addHTML {
		return translatedURL + ".html"
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
	subDomains := map[string]string{
		"en-US": "en.",
		"es-US": "es.",
		"de-DE": "de.",
	}

	for _, sub := range subDomains {
		domain = strings.TrimPrefix(domain, sub)
	}

	sub, ok := subDomains[lang]
	if !ok {
		return "www." + domain
	}
	return sub + domain
}
