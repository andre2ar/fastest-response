package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type BrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a CEP as an argument")
	}

	cep := os.Args[1]

	brasilApiChannel := make(chan BrasilApiResponse)
	viaCepChannel := make(chan ViaCepResponse)

	go GetBrasilApi(cep, brasilApiChannel)
	go GetViaCep(cep, viaCepChannel)

	select {
	case response := <-brasilApiChannel:
		fmt.Println("Brasil API response:")
		fmt.Println(response)
	case response := <-viaCepChannel:
		fmt.Println("Via CEP response:")
		fmt.Println(response)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}

func GetBrasilApi(cep string, channel chan<- BrasilApiResponse) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	res, err := GetRequest(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	responseBody, _ := io.ReadAll(res.Body)
	var brasilApiResponse BrasilApiResponse
	err = json.Unmarshal(responseBody, &brasilApiResponse)
	if err != nil {
		log.Println(err)
	}

	channel <- brasilApiResponse
}

func GetViaCep(cep string, channel chan<- ViaCepResponse) {
	url := "https://viacep.com.br/ws/" + cep + "/json"
	res, err := GetRequest(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	responseBody, _ := io.ReadAll(res.Body)
	var viaCepResponse ViaCepResponse
	err = json.Unmarshal(responseBody, &viaCepResponse)
	if err != nil {
		log.Println(err)
	}

	channel <- viaCepResponse
}

func GetRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, nil
}
