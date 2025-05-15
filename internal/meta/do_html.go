package meta

type DoHtml struct {
	Name    string
	Struct  string
	InArgs  []*Arg
	OutArgs []*Arg
	Comment string
	Q       string
	M       string
	Version string
}
