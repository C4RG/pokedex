package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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
	pokedex    Pokedex
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

type Pokemon struct {
	Name    string `json:"name"`
	BaseExp int    `json:"base_experience"`
}

type PokemonInspect struct {
	Name   string         `json:"name"`
	Height int            `json:"height"`
	Weight int            `json:"weight"`
	Stats  []PokemonStats `json:"stats"`
	Types  []PokemonTypes `json:"types"`
}

type PokemonStats struct {
	BaseStat int         `json:"base_stat"`
	Stat     PokemonStat `json:"stat"`
}

type PokemonStat struct {
	Name string `json:"name"`
}

type PokemonTypes struct {
	Type PokemonType `json:"type"`
}

type PokemonType struct {
	Name string `json:"name"`
}

type Pokedex struct {
	Entries map[string]Pokemon
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Second)
	startingUrl := "https://pokeapi.co/api/v2/location-area/"
	reg := config{commands: getCommands(), next: &startingUrl, cache: cache, pokedex: Pokedex{make(map[string]Pokemon)}}
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
			} else {
				err := command.callback(&reg, "")

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
		fmt.Printf("Exploring from cache%s\n", e)
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

func commandCatch(c *config, e string) error {
	startingUrl := "https://pokeapi.co/api/v2/pokemon/"
	pokemon := Pokemon{}
	fmt.Printf("Throwing a Pokeball at %s...\n", e)
	res_c, exists := c.cache.Get(startingUrl + e)
	if exists {
		err := json.Unmarshal(res_c, &pokemon)
		if err != nil {
			return err
		}
	}
	res, err := http.Get(startingUrl + e)
	if err != nil {
		return err
	}
	if res.Status != "200 OK" {
		fmt.Printf("Pokemon %s was not found!", e)
		return nil
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &pokemon)
	rng := rand.Intn(pokemon.BaseExp)
	if rng > pokemon.BaseExp/5 {
		fmt.Printf("%s was caught!\n", e)
		c.cache.Add(startingUrl+e, body)
		c.pokedex.Entries[e] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", e)
	}

	return nil
}

func commandInspect(c *config, e string) error {
	pokemon := PokemonInspect{}
	url := "https://pokeapi.co/api/v2/pokemon/"
	res, exists := c.cache.Get(url + e)
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	err := json.Unmarshal(res, &pokemon)
	if err != nil {
		return err
	}

	fmt.Printf("Name:%s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Width: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, i := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", i.Stat.Name, i.BaseStat)
	}
	fmt.Println("Types:")
	for _, i := range pokemon.Types {
		fmt.Printf("  - %s\n", i.Type.Name)
	}

	return nil
}
