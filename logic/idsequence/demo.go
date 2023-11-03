package idsequence

import (
	"context"
	"fmt"
)

var IdSequenceMap = make(map[string]*IdSequence)

func Init() {
	IdSequenceMap = map[string]*IdSequence{
		"demo": NewIdSequence(10000, "demo"),
	}
}

func GetIdSequence(biz string) (*IdSequence, error) {
	idSequence, ok := IdSequenceMap[biz]
	if !ok {
		return nil, fmt.Errorf("biz=(%s) nor support", biz)
	}

	return idSequence, nil
}

func Stop() {
	for _, idSequence := range IdSequenceMap {
		idSequence.stopMonitor <- true
		idSequence.Close()
		idSequence.saveLastId(context.Background())
	}
}
