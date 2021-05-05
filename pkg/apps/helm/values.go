package helm

import "strings"

var CliValues []string

func cliValuesMap() map[string]string {
	kvPairs := make(map[string]string, len(CliValues))
	for _, v := range CliValues {
		s := strings.Split(v, "=")
		kvPairs[s[0]] = s[1]
	}
	return kvPairs
}
