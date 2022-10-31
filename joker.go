package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Joke() string {
	url := "https://v2.jokeapi.dev/joke/Programming?type=single"
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var result map[string]interface{}
			json.Unmarshal([]byte(body), &result)

			if result["error"] == false {
				reply := fmt.Sprintf("%v", result["joke"])
				return reply
			} else {
				fmt.Println(result["error"])
			}
		} else {
			log.Fatal(err)
			return "Some error Happened"
		}
	}
	return "Whoooops! You should not be seeing this."
}
