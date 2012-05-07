package schema

import (
    "strings"
)

// copy from the json package

type tagOptions string

// parse a struct field tag and options
func parseTag(tag string) (string, tagOptions) {
    if idx := strings.Index(tag, ","); idx != -1 {
        return tag[:idx], tagOptions(tag[idx+1:])
    }

    return tag, tagOptions("")
}

// Contains returns whether checks that a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
    if len(o) == 0 {
        return false
    }
    s := string(o)
    for s != "" {
        var next string
        i := strings.Index(s, ",")
        if i >= 0 {
            s, next = s[:i], s[i+1:]
        }
        if s == optionName {
            return true
        }
        s = next
    }
    return false
}

