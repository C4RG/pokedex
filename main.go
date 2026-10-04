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
	callback    func(*config, string) error
}

type config struct {
	commands   map[string]cliCommand
	next, prev *string
	cache      pokecache.Cache
}

type locationAreaRes struct {
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []location `json:"results"`
}

type location struct {
	Name       string             `json:"name"`
	PokemonEnc []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon PokemonInfo `json:"pokemon"`
}
type PokemonInfo struct {
	Name string `json:"name"`
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
		cmd := cleanInput(scanner.Text())
		command, exists := commands[cmd[0]]
		if !exists {
			fmt.Println("Unknown command")
		} else {
			if len(cmd) > 1 {
				err := command.callback(&reg, cmd[1])
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			}
			err := command.callback(&reg, "")

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
		"explore": {
			name:        "explore",
			description: "Explore the pokemons in the area specified",
			callback:    commandExplore,
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

func commandMap(c *config, e string) error {
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

func commandMapb(c *config, e string) error {
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

func commandExplore(c *config, e string) error {
	startingUrl := "https://pokeapi.co/api/v2/location-area/"
	location := location{}
	loc, exists := c.cache.Get(startingUrl + e)
	if exists {
		fmt.Printf("Exploring %s\n", e)
		fmt.Printf("Found Pokemon:\n")
		err := json.Unmarshal(loc, &location)
		if err != nil {
			return err
		}
		for _, p := range location.PokemonEnc {
			fmt.Println(p.Pokemon.Name)
		}
		return nil
	}
	res, err := http.Get(startingUrl + e)
	if err != nil {
		return err
	}
	if res.Status != "200 OK" {
		fmt.Printf("Invalid location! %s. Request status: %s", e, res.Status)
		return fmt.Errorf("%s", res.Status)
	}
	fmt.Printf("Exploring %s\n", e)
	fmt.Printf("Found Pokemon:\n")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return err
	}
	for _, p := range location.PokemonEnc {
		fmt.Println(p.Pokemon.Name)
	}
	return nil

}
