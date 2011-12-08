package proto

type Integer int64
type Bulk struct {
	Nil  bool   // Nil Bulk
	Bulk []byte // Bulk content
}
type MultiBulk struct {
	Nil       bool   // Nil Multi-Bulk
	MultiBulk []Bulk // Multi-Bulk content
}
type Status string
type Error string

func (i Integer) String() string {
	return ":" + itoa(int64(i)) + "\r\n"
}

func (b Bulk) String() string {
	if b.Nil {
		return "$-1\r\n"
	}
	return "$" + itoa(int64(len(b.Bulk))) + "\r\n" +
		string(b.Bulk) + "\r\n"
}

func (m MultiBulk) String() string {
	if m.Nil {
		return "*-1\r\n"
	}
	if len(m.MultiBulk) == 0 {
		return "*0"
	}
	s := "*" + itoa(int64(len(m.MultiBulk))) + "\r\n"
	for _, bulk := range m.MultiBulk {
		s += bulk.String()
	}
	return s
}

func (s Status) String() string {
	return "+" + string(s) + "\r\n"
}

func (e Error) String() string {
	return "-" + string(e) + "\r\n"
}
