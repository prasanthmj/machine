package machine

type EmptyLog struct {
}

func NewEmptyLog() *EmptyLog {

	return &EmptyLog{}
}
func (*EmptyLog) Print(...interface{}) {
}

func (*EmptyLog) Printf(string, ...interface{}) {
}

func (*EmptyLog) Println(...interface{}) {
}

func (*EmptyLog) Fatal(...interface{}) {
}

func (*EmptyLog) Fatalf(string, ...interface{}) {
}

func (*EmptyLog) Fatalln(...interface{}) {
}

func (*EmptyLog) Panic(...interface{}) {
}

func (*EmptyLog) Panicf(string, ...interface{}) {
}

func (*EmptyLog) Panicln(...interface{}) {
}
