package request

import (
	"io/ioutil"
	"log"
	"net/http"
)

// User agent is required otherwise the request cannot be made
const UserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:91.0) Gecko/20100101 Firefox/91.0"

// Added function in an attempt to clean up the code a little bit.
func UserAgentGetReq(endpoint string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

// to be fair this method might work for all other requests.
// dominos is really strict with the requests that are made
func DominosGet(endpoint string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-GB,en;q=0.5")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

// curl example
// curl 'https://www.dominos.co.uk/Deals/StoreDealGroups?dealsVersion=637670407835700000&fulfilmentMethod=1&isoCode=en-GB&storeId=28131&v=97.1.0.4'
// -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:91.0) Gecko/20100101 Firefox/91.0'
// -H 'Accept: application/json, text/plain, */*'
// -H 'Accept-Language: en-GB,en;q=0.5'
// --compressed

// curl 'https://www.dominos.co.uk/storefindermap/storesearch?SearchText=me46ea'
// -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:91.0) Gecko/20100101 Firefox/91.0'
// -H 'Accept: application/json, text/plain, */*'
// -H 'Accept-Language: en-GB,en;q=0.5'
// --compressed
