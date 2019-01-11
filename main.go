package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	input "github.com/tcnksm/go-input"
	editorconfig "gopkg.in/editorconfig/editorconfig-core-go.v1"
)

func itIsInt(value string) error {
	_, err := strconv.Atoi(value)

	if err != nil {
		return input.ErrNotNumber
	}

	return nil

}

func main() {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	file := ".editorconfig"
	defaultConfig := &editorconfig.Editorconfig{
		Root: true,
		Definitions: []*editorconfig.Definition{
			&editorconfig.Definition{
				Selector:               "*",
				Charset:                editorconfig.CharsetUTF8,
				IndentStyle:            editorconfig.IndentStyleSpaces,
				IndentSize:             "4",
				TabWidth:               4,
				EndOfLine:              editorconfig.EndOfLineLf,
				TrimTrailingWhitespace: true,
				InsertFinalNewline:     true,
			},
		},
	}

	config, err := editorconfig.ParseFile(file)

	creating := false

	if err != nil {
		if os.IsNotExist(err) {
			creating = true
			config = defaultConfig
		} else {
			log.Printf("Unable to read editorconfig: %v\n", err)
		}
	}

	indentStyle, err := ui.Select("Indent style", []string{editorconfig.IndentStyleSpaces, editorconfig.IndentStyleTab},
		&input.Options{
			Default: config.Definitions[0].IndentStyle,
			Loop:    true,
		})

	indentSize := config.Definitions[0].IndentSize
	tabWidthStr := strconv.Itoa(config.Definitions[0].TabWidth)

	if indentStyle == editorconfig.IndentStyleSpaces {
		indentSize, err = ui.Ask("Indent size", &input.Options{
			Default:      indentSize,
			ValidateFunc: itIsInt,
			Loop:         true,
		})
	} else {

		tabWidthStr, err = ui.Ask("Tab width", &input.Options{
			Default:      tabWidthStr,
			ValidateFunc: itIsInt,
			Loop:         true,
		})
	}

	endOfLine, err := ui.Select("End of line", []string{editorconfig.EndOfLineCr, editorconfig.EndOfLineCrLf, editorconfig.EndOfLineLf},
		&input.Options{
			Default: config.Definitions[0].EndOfLine,
			Loop:    true,
		})

	trimTrailingWhitespaceStr, err := ui.Select("Trim training whitespace", []string{"true", "false"},
		&input.Options{
			Default: strconv.FormatBool(config.Definitions[0].TrimTrailingWhitespace),
			Loop:    true,
		})

	insertFinalNewlineStr, err := ui.Select("Insert final new line", []string{"true", "false"},
		&input.Options{
			Default: strconv.FormatBool(config.Definitions[0].InsertFinalNewline),
			Loop:    true,
		})

	str := fmt.Sprintf(
		"root = %t\n"+
			"\n"+
			"[*]\n"+
			"charset = %s\n"+
			"indent_style = %s\n"+
			"indent_size = %s\n"+
			"tab_width = %s\n"+
			"end_of_line = %s\n"+
			"insert_final_newline = %s\n"+
			"trim_trailing_whitespace = %s\n",
		config.Root,
		config.Definitions[0].Charset,
		indentStyle,
		indentSize,
		tabWidthStr,
		endOfLine,
		insertFinalNewlineStr,
		trimTrailingWhitespaceStr)

	err = ioutil.WriteFile(file, []byte(str), 0644)

	if err != nil {
		log.Panicf("Cannot save file: %v", err)
	}

	if creating {
		fmt.Println("Created editorconfig in", file)
	} else {
		fmt.Println("Updated editorconfig in", file)
	}
}
