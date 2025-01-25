package helpers

import (
	"fmt"
	"io"
	"net/http"
)

func SendGetQuery(url string) (string, error) {
	const fn = "external.sendQuery"

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("%s error send query %s", fn, err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s query return status code: %d", fn, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("%s error read body %s", fn, err.Error())
	}

	return string(body), nil
}
