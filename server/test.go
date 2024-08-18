package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	baseURL := "http://localhost:3000/api"

	// Make the /api/maze call
	mazeURL := fmt.Sprintf("%s/maze?mazeSize=3&singlePath=true", baseURL)
	resp, err := http.Get(mazeURL)
	if err != nil {
		fmt.Printf("Error making /api/maze request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d from /api/maze\n", resp.StatusCode)
		return
	}

	//mazeBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading /api/maze response body: %v\n", err)
		return
	}
	/*
		fmt.Println("Response from /api/maze:")
		fmt.Println(string(mazeBody))
	*/

	// Make the /api/solution call
	solutionURL := fmt.Sprintf("%s/solution", baseURL)
	resp, err = http.Get(solutionURL)
	if err != nil {
		fmt.Printf("Error making /api/solution request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d from /api/solution\n", resp.StatusCode)
		return
	}

	solutionBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading /api/solution response body: %v\n", err)
		return
	}

	fmt.Println("Response from /api/solution:")
	fmt.Println(string(solutionBody))

}
