package subcmd

import "flag"

type ValueType interface {
	flag.Value
	Copy() ValueType
}
