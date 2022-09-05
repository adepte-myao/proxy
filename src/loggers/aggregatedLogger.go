package loggers

import "log"

type AggregatedLoggers struct {
	lgs []*log.Logger
}

func NewAggregatedLoggers(l ...*log.Logger) *AggregatedLoggers {
	return &AggregatedLoggers{l}
}

func (ls *AggregatedLoggers) Println(v ...any) {
	for i := 0; i < len(ls.lgs); i++ {
		ls.lgs[i].Println(v...)
	}
}

func (ls *AggregatedLoggers) Printf(format string, v ...any) {
	for i := 0; i < len(ls.lgs); i++ {
		ls.lgs[i].Printf(format, v...)
	}
}
