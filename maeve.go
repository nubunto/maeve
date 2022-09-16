package maeve

import (
	"fmt"
	"strings"
)

var DefaultSeparator = "/"

type StringPath string

func KV(keysAndValues ...string) KeyValueList {
	if len(keysAndValues)%2 != 0 {
		panic(fmt.Sprintf("invalid number of arguments in call to KV: %d", len(keysAndValues)))
	}

	kv := make(KeyValueList, 0, len(keysAndValues)/2)

	for i := 0; i < len(keysAndValues); i += 2 {
		key := keysAndValues[i]
		value := keysAndValues[i+1]

		kv = append(kv, KeyValue{
			Path:  key,
			Value: value,
		})
	}

	return kv
}

func Path(paths ...string) StringPath {
	return StringPath(strings.Join(paths, DefaultSeparator))
}

func TrimDynamic(path StringPath) StringPath {
	return StringPath(strings.Replace(string(path), "*", "", 1))
}

func IsDynamic(path StringPath) bool {
	return strings.Contains(string(path), "*")
}

type KeyValue struct {
	Path  string
	Value string
}

type KeyValueList []KeyValue
