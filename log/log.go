package log

import (
	"cosmossdk.io/log"
)

type Logger = log.Logger

type Fields map[string]any

func (f Fields) With(xs ...any) Fields {
	f2 := make(Fields, len(f)+(len(xs)/2))
	for k, v := range f {
		f2[k] = v
	}
	for i := 0; i < len(xs)/2; i++ {
		key, is := xs[i*2].(string)
		if !is {
			continue
		}
		val := xs[i*2+1]
		f2[key] = val
	}
	return f2
}

func (f Fields) Merge(f2 Fields) Fields {
	f3 := make(Fields, len(f)+len(f2))
	for k, v := range f {
		f3[k] = v
	}
	for k, v := range f2 {
		f3[k] = v
	}
	return f3
}

func (f Fields) Slice() []any {
	s := make([]any, len(f)*2)
	i := 0
	for k, v := range f {
		s[i] = k
		s[i+1] = v
		i += 2
	}
	return s
}
