package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func InputLoop(scanner bufio.Scanner, reg *config) {
	for {
		commands := reg.commands
		fmt.Print("\nPokedex > ")
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			fmt.Println(err)
		}
		cmd := cleanInput(scanner.Text())
		command, exists := commands[cmd[0]]
		if !exists {
			fmt.Println("Unknown command")
		} else {
			if len(cmd) > 1 {
				err := command.callback(reg, cmd[1])
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			} else {
				err := command.callback(reg, "")

				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			}

		}
	}
}

func cleanInput(text string) []string {
	ltext := strings.ToLower(text)
	split := strings.Fields(ltext)
	return split
}
