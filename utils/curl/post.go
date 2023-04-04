package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Post(url string, body interface{}, header http.Header, responseBody any) int {
	client := &http.Client{}

	b, e := json.Marshal(body)

	if e != nil {
		fmt.Println("err")
	}

	bReader := bytes.NewReader(b)

	res, err := http.NewRequest("POST", url, bReader)

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
