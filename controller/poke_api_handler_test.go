package controller

import (
	"catching-pokemons/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

// func TestGetPokemonFromPokeApiSuccess(t *testing.T) {
// 	c := require.New(t)

// 	pokemon, err := GetPokemonFromPokeApi("pikachu")
// 	c.NoError(err)

// 	body, err := ioutil.ReadFile("samples/pokeapi_read.json")
// 	c.NoError(err)

// 	var expected models.PokeApiPokemonResponse
// 	err = json.Unmarshal([]byte(body), &expected)
// 	c.NoError(err)

// 	c.Equal(expected, pokemon)
// }

func TestGetPokemonFromApiSuccess(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "pikachu"

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	c.NoError(err)

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder("GET", request, httpmock.NewStringResponder(200, string(body)))

	pokemon, err := GetPokemonFromPokeApi(id)
	c.NoError(err)

	var expected models.PokeApiPokemonResponse
	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(expected, pokemon)
}

func TestGetPokemonFromApiInternalServerError(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "pikachu"

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	c.NoError(err)

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder("GET", request, httpmock.NewStringResponder(500, string(body)))

	_, err = GetPokemonFromPokeApi(id)
	c.NotNil(err)
	c.EqualError(ErrPokeApiFailure, err.Error())
}

func TestGetPokemonFromApiNotFound(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "pikachu"

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	c.NoError(err)

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder("GET", request, httpmock.NewStringResponder(404, string(body)))

	_, err = GetPokemonFromPokeApi(id)
	c.NotNil(err)
	c.EqualError(ErrPokemonNotFound, err.Error())
}

func TestGetPokemon(t *testing.T) {
	c := require.New(t)

	r, err := http.NewRequest("GET", "/pokemon{id}", nil)
	c.NoError(err)

	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": "pikachu",
	}

	r = mux.SetURLVars(r, vars)

	GetPokemon(w, r)

	body, err := ioutil.ReadFile("samples/api_response.json")
	c.NoError(err)

	var expectedPokemon models.Pokemon
	err = json.Unmarshal([]byte(body), &expectedPokemon)
	c.NoError(err)

	var actualPokemon models.Pokemon
	err = json.Unmarshal(w.Body.Bytes(), &actualPokemon)
	c.NoError(err)

	c.Equal(http.StatusOK, w.Code)
	c.Equal(expectedPokemon, actualPokemon)
}
