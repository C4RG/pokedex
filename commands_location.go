package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
