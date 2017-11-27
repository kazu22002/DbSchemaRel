package singular

import (
	"regexp"
)

/**
 * ループで回した時に、ランダムに抽出されたためキーのみ配列を作成する。
 */
var singular_rules = map[string]string{
	"(s)tatuses$":     "${1}tatus",
	"^(.*)(menu)s$":   "${1}",
	"(quiz)zes$":      "${1}",
	"(matr)ices$":     "${1}ix",
	"(vert|ind)ices$": "${1}ex",
	"^(ox)en":         "${1}",
	"(alias)(es)*$":   "${1}",
	"(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$": "${1}us",
	"([ftw]ax)es":        "${1}",
	"(cris|ax|test)es$":  "${1}is",
	"(shoe|slave)s$":     "${1}",
	"(o)es$":             "${1}",
	"ouses$":             "ouse",
	"([^a])uses$":        "${1}us",
	"([m|l])ice$":        "${1}ouse",
	"(x|ch|ss|sh)es$":    "${1}",
	"(m)ovies$":          "${1}ovie",
	"(s)eries$":          "${1}eries",
	"([^aeiouy]|qu)ies$": "${1}y",
	"([lr])ves$":         "${1}f",
	"(tive)s$":           "${1}",
	"(hive)s$":           "${1}",
	"(drive)s$":          "${1}",
	"([^fo])ves$":        "${1}fe",
	"(^analy)ses$":       "${1}sis",
	"(analy|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$": "${1}sis",
	"([ti])a$":    "${1}um",
	"(p)eople$":   "${1}erson",
	"(m)en$":      "${1}an",
	"(c)hildren$": "${1}hild",
	"(n)ews$":     "${1}ews",
	"eaus$":       "eau",
	"^(.*us)$":    "${1}",
	"s$":          ""}

var singular_rules_sort = []string{
	"(s)tatuses$",
	"^(.*)(menu)s$",
	"(quiz)zes$",
	"(matr)ices$",
	"(vert|ind)ices$",
	"^(ox)en",
	"(alias)(es)*$",
	"(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$",
	"([ftw]ax)es",
	"(cris|ax|test)es$",
	"(shoe|slave)s$",
	"(o)es$",
	"ouses$",
	"([^a])uses$",
	"([m|l])ice$",
	"(x|ch|ss|sh)es$",
	"(m)ovies$",
	"(s)eries$",
	"([^aeiouy]|qu)ies$",
	"([lr])ves$",
	"(tive)s$",
	"(hive)s$",
	"(drive)s$",
	"([^fo])ves$",
	"(^analy)ses$",
	"(analy|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$",
	"([ti])a$",
	"(p)eople$",
	"(m)en$",
	"(c)hildren$",
	"(n)ews$",
	"eaus$",
	"^(.*us)$",
	"s$",
}

func SingleName(name string) string {
	var single_name = name

	for _, key := range singular_rules_sort {
		if regexp.MustCompile(key).MatchString(name) {
			single_name = regexp.MustCompile(key).ReplaceAllString(name, singular_rules[key])
			break
		}
	}

	return single_name
}
