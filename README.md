# Pokedex

A command-line Pokedex written in Go. It uses the [PokeAPI](https://pokeapi.co/) to browse location areas and catch Pokémon.

## Run

Requires Go 1.25 or later and an internet connection.

From the project directory, start the REPL with:

```sh
go run .
```

Enter commands at the `Pokedex >` prompt. Commands are case-insensitive.

## Commands

| Command | Description |
| --- | --- |
| `help` | Show the available commands. |
| `map` | Show 20 location areas; run again to move to the next page. |
| `mapb` | Show the previous 20 location areas, when available. |
| `explore <area>` | List Pokémon found in a location area, for example `explore canalave-city-area`. |
| `catch <pokemon>` | Try to catch a Pokémon, for example `catch pikachu`. Catching is chance-based. |
| `inspect <pokemon>` | Show details for a Pokémon you have caught. |
| `pokedex` | List the Pokémon you have caught. |
| `exit` | Quit the REPL. |

## Cache

The application keeps cached API responses in memory while it is running. Entries older than five seconds are removed, so the cache is temporary and is not saved between runs. Location-area responses are cached for map navigation, and a successfully caught Pokémon's response is cached for inspection.