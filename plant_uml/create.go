package plant_uml

import (
	"../schema"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var singular_rules = map[string]string{
	"(s)tatuses$":     "12tatus",
	"^(.*)(menu)s$":   "12",
	"(quiz)zes$":      "1",
	"(matr)ices$":     "1ix",
	"(vert|ind)ices$": "1ex",
	"^(ox)en":         "1",
	"(alias)(es)*$":   "1",
	"(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$": "1us",
	"([ftw]ax)es":        "1",
	"(cris|ax|test)es$":  "1is",
	"(shoe|slave)s$":     "1",
	"(o)es$":             "1",
	"ouses$":             "ouse",
	"([^a])uses$":        "1us",
	"([m|l])ice$":        "1ouse",
	"(x|ch|ss|sh)es$":    "1",
	"(m)ovies$":          "12ovie",
	"(s)eries$":          "12eries",
	"([^aeiouy]|qu)ies$": "1y",
	"([lr])ves$":         "1f",
	"(tive)s$":           "1",
	"(hive)s$":           "1",
	"(drive)s$":          "1",
	"([^fo])ves$":        "1fe",
	"(^analy)ses$":       "1sis",
	"(analy|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$": "12sis",
	"([ti])a$":    "1um",
	"(p)eople$":   "12erson",
	"(m)en$":      "1an",
	"(c)hildren$": "12hild",
	"(n)ews$":     "12ews",
	"eaus$":       "eau",
	"^(.*us)$":    "1",
	"s$":          ""}

func singleName(name string) string {
	var single_name = name

	for key, replace := range singular_rules {
		if regexp.MustCompile(key).MatchString(name) {
			single_name = regexp.MustCompile(key).ReplaceAllString(name, replace)
			break
		}
	}

	return single_name
}

func Output(data []schema.DbSchema) {

	fmt.Println(data)

	if len(data) < 0 {
		return
	}
	var content bytes.Buffer
	recode := "\r\n"
	count := len(data)
	for i := 0; i < count; i++ {
		//		content.WriteString(data[i].table)
		content.WriteString(recode)
	}
	ioutil.WriteFile("planet_uml.txt", []byte(content.String()), os.ModePerm)

	// single_name := singleName(name)
}
