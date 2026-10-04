package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	pokecache "c4rg/pokedex/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	commands   map[string]cliCommand
	next, prev *string
	cache      pokecache.Cache
}

type locationAreaRes struct {
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []location `json:results"`
}

type location struct {
	Name string `json:name`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Second)
	startingUrl := "https://pokeapi.co/api/v2/location-area/"
	reg := config{commands: getCommands(), next: &startingUrl, cache: cache}
	commands := reg.commands

	for {
		fmt.Print("\nPokedex > ")
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			fmt.Println(err)
		}
		cmd := scanner.Text()
		command, exists := commands[cmd]
		if !exists {
			fmt.Println("Unknown command")
		} else {
			err := command.callback(&reg)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}
	}
}

func cleanInput(text string) []string {
	ltext := strings.ToLower(text)
	split := strings.Fields(ltext)
	return split
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
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
	}
}

func commandExit(c *config) error {
	myerr := fmt.Errorf("Closing the Pokedex... Goodbye!")
	return myerr
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := c.commands
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(c *config) error {
	data := locationAreaRes{}
	next, exists := c.cache.Get(*c.next)
	if exists {
		err := json.Unmarshal(next, &data)
		if err != nil {
			return err
		}
		c.next = data.Next
		c.prev = data.Previous
		for _, l := range data.Results {
			fmt.Println(l.Name)
		}
		return nil
	}
	res, err := http.Get(*c.next)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	c.cache.Add(*c.next, body)
	c.next = data.Next
	c.prev = data.Previous
	for _, l := range data.Results {
		fmt.Println(l.Name)
	}
	return nil
}

func commandMapb(c *config) error {
	data := locationAreaRes{}
	if c.prev == nil {
		fmt.Println("You're on the first page")
		return nil
	}
	prev, exists := c.cache.Get(*c.prev)
	if exists {
		err := json.Unmarshal(prev, &data)
		if err != nil {
			return err
		}
		for _, l := range data.Results {
			fmt.Println(l.Name)
		}
		return nil
	}
	res, err := http.Get(*c.prev)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	c.cache.Add(*c.prev, body)
	c.next = data.Next
	c.prev = data.Previous
	for _, l := range data.Results {
		fmt.Println(l.Name)
	}
	return nil
}
