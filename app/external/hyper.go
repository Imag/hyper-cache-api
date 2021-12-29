package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Imag/hyper-cache-api/app/types"
)

func RetrieveLicense(key string) (types.License, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.hyper.co/v6/licenses/%s", key), nil)
	if err != nil {
		return types.License{}, err
	}
	
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("HYPER_SECRET_KEY")))

	res, err := client.Do(req)
	if err != nil {
		return types.License{}, err
	}

	var license types.License
	json.NewDecoder(res.Body).Decode(&license)
	
	return license, nil
}