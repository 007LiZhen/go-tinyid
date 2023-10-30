package idsequence

var IdSequenceMap = make(map[string]*IdSequence)

func Init() {
	IdSequenceMap = map[string]*IdSequence{
		"demo": NewIdSequence(10000, "demo"),
	}
}
