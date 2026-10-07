package main

import (
	pokecache "c4rg/pokedex/internal"
	"fmt"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	commands   map[string]cliCommand
	next, prev *string
	cache      pokecache.Cache
	pokedex    Pokedex
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Shows you 20 location areas. Re-run to get 20 more.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows you the 20 previous locations, if there are any.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore the pokemons in the area specified",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Throw pokeball at pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows registered pokemon",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
	}
}

func commandExit(c *config, e string) error {
	myerr := fmt.Errorf("Closing the Pokedex... Goodbye!")
	return myerr
}

func commandHelp(c *config, e string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := c.commands
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}
