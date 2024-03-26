package main

import (
	"html/template"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	enUS = "en-US"
	esUS = "es-US"
	deDE = "de-DE"
)

func fetchLang(url string) string {
	if strings.HasPrefix(url, "es.") {
		return esUS
	} else if strings.HasPrefix(url, "de.") {
		return deDE
	}

	return enUS
}

func lastOne(index int, size int) bool {
	return index == size-1
}

func translate(lang, en, es, de string) string {
	switch lang {
	case esUS:
		return es
	case deDE:
		return de
	default:
		return en
	}
}

func translateHTML(lang, en, es, de string) template.HTML {
	switch lang {
	case esUS:
		return template.HTML(es)
	case deDE:
		return template.HTML(de)
	default:
		return template.HTML(en)
	}
}

var dictionary = buildDictionary([]map[string]string{
	//////////////// PAGES ///////////////////////
	{
		enUS: "home",
		esUS: "inicio",
		deDE: "start",
	}, {
		enUS: "about",
		esUS: "conóceme", // non ASCII
		deDE: "über",     // non ASCII
	}, {
		enUS: "posts",
		esUS: "entradas",
		deDE: "posten",
	}, {
		enUS: "friends",
		esUS: "amigos",
		deDE: "freunde",
	}, {
		enUS: "library",
		esUS: "biblioteca",
		deDE: "bibliothek",
	}, {
		enUS: "recommend",
		esUS: "recomendar",
		deDE: "empfehlen",
	}, {
		enUS: "tags",
		esUS: "etiquetas",
		deDE: "stichwörter", // non ASCII
	}, {
		enUS: "contact",
		esUS: "contacto",
		deDE: "kontakt",
	}, {
		enUS: "todo",
		esUS: "pendiente",
		deDE: "aufgaben",
		//////////////// MISC ///////////////////////
	}, {
		enUS: "hi",
		esUS: "hola",
		deDE: "hallo",
		//////////////// TAGS ///////////////////////
	}, { // medium
		enUS: "articles",
		esUS: "artículos", // non ASCII
		deDE: "artikel",
	}, {
		enUS: "photos",
		esUS: "fotos",
		deDE: "fotos",
	}, { // lang
		enUS: "english",
		esUS: "inglés", // non ASCII
		deDE: "englisch",
	}, {
		enUS: "spanish",
		esUS: "español", // non ASCII
		deDE: "spanisch",
	}, {
		enUS: "german",
		esUS: "alemán", // non ASCII
		deDE: "deutsch",
	}, { // tags
		enUS: "math",
		esUS: "matemáticas", // non ASCII
		deDE: "mathe",
	}, {
		enUS: "milwaukee",
		esUS: "milwaukee",
		deDE: "milwaukee",
	}, {
		enUS: "history",
		esUS: "historia",
		deDE: "geschichte",
	}, {
		enUS: "technology",
		esUS: "tecnologia",
		deDE: "technologie",
	}, {
		enUS: "code",
		esUS: "código", // non ASCII
		deDE: "code",
	}, {
		enUS: "updates",
		esUS: "actualizaciones",
		deDE: "aktualisierungen",
	}, {
		enUS: "personal",
		esUS: "personal",
		deDE: "persönliches", // non ASCII
		///////////////// TIME ////////////////////
	}, { // days of week
		enUS: "monday",
		esUS: "lunes",
		deDE: "montag",
	}, {
		enUS: "tuesday",
		esUS: "martes",
		deDE: "dienstag",
	}, {
		enUS: "wednesday",
		esUS: "miércoles", // non ASCII
		deDE: "mittwoch",
	}, {
		enUS: "thursday",
		esUS: "jueves",
		deDE: "donnerstag",
	}, {
		enUS: "friday",
		esUS: "viernes",
		deDE: "freitag",
	}, {
		enUS: "saturday",
		esUS: "sábado", // non ASCII
		deDE: "samstag",
	}, {
		enUS: "sunday",
		esUS: "domingo",
		deDE: "sonntag",
	}, { // months
		enUS: "january",
		esUS: "enero",
		deDE: "januar",
	}, {
		enUS: "february",
		esUS: "febrero",
		deDE: "februar",
	}, {
		enUS: "march",
		esUS: "marzo",
		deDE: "märz", // non ASCII
	}, {
		enUS: "april",
		esUS: "abril",
		deDE: "april",
	}, {
		enUS: "may",
		esUS: "mayo",
		deDE: "mai",
	}, {
		enUS: "june",
		esUS: "junio",
		deDE: "juni",
	}, {
		enUS: "july",
		esUS: "julio",
		deDE: "juli",
	}, {
		enUS: "august",
		esUS: "agosto",
		deDE: "august",
	}, {
		enUS: "september",
		esUS: "septiembre",
		deDE: "september",
	}, {
		enUS: "october",
		esUS: "octubre",
		deDE: "oktober",
	}, {
		enUS: "november",
		esUS: "noviembre",
		deDE: "november",
	}, {
		enUS: "december",
		esUS: "diciembre",
		deDE: "dezember",
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
	if strings.HasSuffix(last, htmlExt) {
		addHTML = true
		urlPath[len(urlPath)-1] = strings.TrimSuffix(last, htmlExt)
	}

	for _, urlSlice := range urlPath[1:] {
		urlSlice = anglicize(translateKeyword(lang, urlSlice))
		translatedURL = translatedURL + "/" + urlSlice
	}

	if addHTML {
		return translatedURL + htmlExt
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
	case esUS:
		formattedDate = date.Format("Monday, 2 de January del 2006")
		diaDeLaSemana := translateKeyword(esUS, dayOfWeek)
		formattedDate = strings.ReplaceAll(formattedDate, dayOfWeek, diaDeLaSemana)
		mes := translateKeyword(esUS, strings.ToLower(month))
		formattedDate = strings.ReplaceAll(formattedDate, month, mes)
		return formattedDate, nil
	case deDE:
		formattedDate = date.Format("Monday, 02.01.2006")
		tagDerWoche := translateKeyword(deDE, dayOfWeek)
		formattedDate = strings.ReplaceAll(formattedDate, dayOfWeek, tagDerWoche)
		return formattedDate, nil
	default:
		return formattedDate, nil
	}
}

func translateHost(lang, domain string) string {
	subDomains := map[string]string{
		enUS: "en.",
		esUS: "es.",
		deDE: "de.",
	}

	domain = strings.TrimPrefix(domain, "www.")
	for _, sub := range subDomains {
		domain = strings.TrimPrefix(domain, sub)
	}

	sub, ok := subDomains[lang]
	if !ok {
		return "www." + domain
	}
	return sub + domain
}
