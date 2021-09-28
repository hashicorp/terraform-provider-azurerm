// Package ipify provides a single function for retrieving your computer's
// public IP address from the ipify service: http://www.ipify.org
package ipify

import (
	"errors"
	"github.com/jpillora/backoff"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// GetIp queries the ipify service (http://www.ipify.org) to retrieve this
// machine's public IP address.  Returns your public IP address as a string, and
// any error encountered.  By default, this function will run using exponential
// backoff -- if this function fails for any reason, the request will be retried
// up to 3 times.
//
// Usage:
//
//		package main
//
//		import (
//			"fmt"
//			"github.com/rdegges/go-ipify"
//		)
//
//		func main() {
//			ip, err := ipify.GetIp()
//			if err != nil {
//				fmt.Println("Couldn't get my IP address:", err)
//			} else {
//				fmt.Println("My IP address is:", ip)
//			}
//		}
func GetIp() (string, error) {
	b := &backoff.Backoff{
		Jitter: true,
	}
	client := &http.Client{}

	req, err := http.NewRequest("GET", API_URI, nil)
	if err != nil {
		return "", errors.New("Received an invalid status code from ipify: 500. The service might be experiencing issues.")
	}

	req.Header.Add("User-Agent", USER_AGENT)

	for tries := 0; tries < MAX_TRIES; tries++ {
		resp, err := client.Do(req)
		if err != nil {
			d := b.Duration()
			time.Sleep(d)
			continue
		}

		defer resp.Body.Close()

		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New("Received an invalid status code from ipify: 500. The service might be experiencing issues.")
		}

		if resp.StatusCode != 200 {
			return "", errors.New("Received an invalid status code from ipify: " + strconv.Itoa(resp.StatusCode) + ". The service might be experiencing issues.")
		}

		return string(ip), nil
	}

	return "", errors.New("The request failed because it wasn't able to reach the ipify service. This is most likely due to a networking error of some sort.")
}
