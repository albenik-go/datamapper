package tag

import (
	"strings"
)

type Option struct {
	Name  string
	Value string
}

type OptionsList []Option

func ParseOptions(tag string) OptionsList {
	list := strings.Split(tag, ",")
	opts := make([]Option, len(list))

	for i, opt := range list {
		val := strings.SplitN(opt, "=", 2)
		opts[i].Name = val[0]
		if len(val) > 1 {
			opts[i].Value = val[1]
		}
	}

	return opts
}

func (o OptionsList) Lookup(name string) *Option {
	for _, opt := range o {
		if opt.Name == name {
			return &opt
		}
	}
	return nil
}
