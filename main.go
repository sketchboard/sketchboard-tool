package main

import (
	"log"
	"os"

	"sketchboard.io/sketchboard-tool/convertcsv"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

func main() {
	var parser = flags.NewNamedParser("sketchboard", flags.Default)

	inits := []func(*flags.Parser) error{
		// bookmark.InitBookmark,
		convertcsv.Init,
	}

	err := addCommands(inits, parser)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			// do not show error twice
			// fmt.Println(err)
			os.Exit(1)
		}
	}

}

func addCommands(commands []func(*flags.Parser) error, parser *flags.Parser) error {
	for _, init := range commands {
		err := init(parser)
		if err != nil {
			return errors.Wrapf(err, "addCommands")
		}
	}

	return nil
}
