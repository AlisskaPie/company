package ipapi

import (
	"io/ioutil"
	"net/http"
)

// IsExpectedLocation checks if location has expected value or not, return local location
func IsExpectedLocation(locationName string) (string, bool) {
	loc, err := Location()
	if err != nil {
		return "", false
	}
	return loc, loc == locationName
}

func Location() (string, error) {
	ipapiClient := http.Client{}
	url := "https://ipapi.co/country_name"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "ipapi.co/#go-v1.5")

	res, err := ipapiClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
