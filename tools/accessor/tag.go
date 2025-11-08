package main

import (
	"errors"
	"strconv"
	"strings"
)

const tagName = "accessor"
const tagInline = "inline"
const tagIgnore = "-"
const tagOptionEqKey = "eq"
const tagOptionEqValue = "=="             // 使用==比较
const tagOptionEqForceValue = "-"         // 不做比较强制更新
const tagOptionEqualMethodValue = "Equal" // 使用Equal方法比较

var (
	errTagSyntax      = errors.New("bad syntax for struct tag pair")
	errTagKeySyntax   = errors.New("bad syntax for struct tag key")
	errTagValueSyntax = errors.New("bad syntax for struct tag value")

	errKeyNotSet      = errors.New("tag key does not exist")
	errTagNotExist    = errors.New("tag does not exist")
	errTagKeyMismatch = errors.New("mismatch between key and tag.key")
)

// Tags represent a set of tags from a single struct field
type Tags struct {
	tags []*Tag
}

// Tag defines a single struct's string literal tag
type Tag struct {
	// Key is the tag key, such as json, xml, etc..
	// i.e: `json:"foo,omitempty". Here key is: "json"
	Key string

	// Name is a part of the value
	// i.e: `json:"foo,omitempty". Here name is: "foo"
	Name string

	// i.e: `json:"foo,omitempty". Here value is: "omitempty"
	Values []string

	// Options is a part of the value. It contains a slice of tag options i.e:
	// `json:"foo,inline;". Here options is: map["inline":"" ]
	Options map[string]string
}

// Parse parses a single struct field tag and returns the set of tags.
func ParseTag(tag string) (*Tags, error) {
	var tags []*Tag

	hasTag := tag != ""

	// NOTE(arslan) following code is from reflect and vet package with some
	// modifications to collect all necessary information and extend it with
	// usable methods
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax
		// error. Strictly speaking, control chars include the range [0x7f,
		// 0x9f], not just [0x00, 0x1f], but in practice, we ignore the
		// multi-byte control characters as it is simpler to inspect the tag's
		// bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}

		if i == 0 {
			return nil, errTagKeySyntax
		}
		if i+1 >= len(tag) || tag[i] != ':' {
			return nil, errTagSyntax
		}
		if tag[i+1] != '"' {
			return nil, errTagValueSyntax
		}

		key := tag[:i]
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			return nil, errTagValueSyntax
		}

		qvalue := tag[:i+1]
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return nil, errTagValueSyntax
		}

		tag := &Tag{Key: key, Values: []string{}, Options: map[string]string{}}

		res := strings.Split(value, ";")
		tag.parseNameValue(res[0])
		tag.parseOptions(res[1:])

		tags = append(tags, tag)
	}

	if hasTag && len(tags) == 0 {
		return nil, nil
	}

	return &Tags{
		tags: tags,
	}, nil
}

// Get returns the tag associated with the given key. If the key is present
// in the tag the value (which may be empty) is returned. Otherwise, the
// returned value will be the empty string. The ok return value reports whether
// the tag exists or not (which the return value is nil).
func (t *Tags) Get(key string) (*Tag, error) {
	for _, tag := range t.tags {
		if tag.Key == key {
			return tag, nil
		}
	}

	return nil, errTagNotExist
}

// Tags returns a slice of tags. The order is the original tag order unless it
// was changed.
func (t *Tags) Tags() []*Tag {
	return t.tags
}

// Keys returns a slice of tags' keys.
func (t *Tags) Keys() []string {
	var keys []string
	for _, tag := range t.tags {
		keys = append(keys, tag.Key)
	}
	return keys
}

// HasOption returns true if the given option is available in options
func (t *Tag) HasValue(v string) bool {
	for _, tagOpt := range t.Values {
		if tagOpt == v {
			return true
		}
	}

	return false
}

// HasOption returns true if the given option is available in options
func (t *Tag) HasOption(opt string) bool {
	for k, _ := range t.Options {
		if k == opt {
			return true
		}
	}

	return false
}

// OptionValue returns option value
func (t *Tag) OptionValue(opt string) string {
	v, ok := t.Options[opt]
	if ok {
		return v
	}

	return ""
}

func (t *Tag) parseNameValue(str string) {
	res := strings.Split(str, ",")
	t.Name = res[0]
	t.Values = res[1:]
}

func (t *Tag) parseOptions(options []string) {
	tags := map[string]string{}
	for _, o := range options {
		op := strings.Split(o, ":")
		if len(op) == 0 {
			continue
		}
		v := ""
		if len(op) > 1 {
			v = op[1]
		}

		tags[op[0]] = v
	}
	t.Options = tags
}
