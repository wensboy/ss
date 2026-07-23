package config

type Source func() (any, bool)

func Lookup(ss ...Source) (any, bool) {
	for _, s := range ss {
		v, ok := s()
		if ok {
			return v, true
		}
	}
	return nil, false
}

func MustLookup(ss ...Source) any {
	v, found := Lookup(ss...)
	if found {
		return v
	}
	panic("config - not found source")
}
