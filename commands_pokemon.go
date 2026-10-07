package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

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
		fmt.Println("You may now inspect it with the inspect command")
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

func commandPokedex(c *config, e string) error {
	if len(c.pokedex.Entries) < 1 {
		fmt.Printf("No pokemons")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, e := range c.pokedex.Entries {
		fmt.Printf(" - %s\n", e.Name)
	}

	return nil
}
