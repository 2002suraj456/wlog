package main

import (
	"fmt"
	"os"
	"wlog/editor"
	"wlog/store"
)

func main() {
	st := store.NewStore(store.Option{
		Directory: "./wlog_data",
	})
	switch processArgs() {
	case Input:
		{
			content, err := editor.Read()
			if err != nil {
				fmt.Println("error reading from editor:", err)
				return
			}

			fmt.Println("Content from editor:")
			fmt.Println(content)

			err = st.SaveData([]byte(content))
			if err != nil {
				fmt.Println("error saving to store:", err)
				return
			}

			fmt.Println("Store saved.")
		}
	case Manage:
		{
			// some tui based manager for logs

		}
	case Unknown:
		{
			fmt.Println("unknown command")
			fmt.Printf(usage)
		}
	}
}

const usage = `
Usage:
  wlog            # Open editor to input log
  wlog manage     # Manage existing logs
`

type Action int

const (
	Input  Action = iota // Normal editor input mode
	Manage               // Manage existing logs
	Unknown
)

func processArgs() Action {
	args := os.Args[1:]
	if len(args) == 0 {
		return Input
	}

	if len(args) == 1 {
		switch args[0] {
		case "manage":
			return Manage
		}
	}

	return Unknown
}
