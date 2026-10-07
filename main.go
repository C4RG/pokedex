package main

import (
	"bufio"
	"os"
	"time"

	pokecache "c4rg/pokedex/internal"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Second)
	startingUrl := "https://pokeapi.co/api/v2/location-area/"
	reg := config{commands: getCommands(), next: &startingUrl, cache: cache, pokedex: Pokedex{make(map[string]Pokemon)}}
	InputLoop(*scanner, &reg)

}
