package currency

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	_ "github.com/go-co-op/gocron"
)

//runClient will send a get request on apiURL
//and will return response body bites
func runClient(apiURL string) ([]byte, error) {
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyBites, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bodyBites, nil
}

//convertClientResponse will convert bodyBites from client response
//into a RatesResponse struct
func convertClientResponse(bodyBites []byte) (RatesResponse, error) {
	var dataResponse RatesResponse
	if err := json.Unmarshal(bodyBites, &dataResponse); err != nil {
		return RatesResponse{}, err
	}
	return dataResponse, nil
}
