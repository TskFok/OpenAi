package curl

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Get(url string, header http.Header, responseBody any) int {
	client := &http.Client{}

	res, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("error")
	}

	res.Header = header

	rep, _ := client.Do(res)

	fmt.Println(rep)
	if rep.StatusCode == http.StatusOK {
		decode := json.NewDecoder(rep.Body)

		decodeErr := decode.Decode(responseBody)

		if decodeErr != nil {
			fmt.Println(decodeErr.Error())
		}
	}
	resErr := rep.Body.Close()

	if resErr != nil {
		fmt.Println(resErr.Error())
	}

	return rep.StatusCode
}
