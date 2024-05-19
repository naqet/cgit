package commands

type Arg struct {
	name         string
	defaultValue string
	choices      []string
	help         string
}

type Processible interface {
	Process([]string) error
}

type Command struct {
	name string
	help string
	args []Arg
}

var Commands = map[string]Processible{
	"init":     &initCmd,
	"cat-file": &catFileCmd,
}
