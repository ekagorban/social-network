package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"social-network/internal/service/signup"

	"github.com/bxcodec/faker/v3"

	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	num := 1000000
	sleep := 5 * time.Millisecond

	for i := 0; i < num; i++ {
		data := formData(i)
		go send(data)
		time.Sleep(sleep)
	}
}

func send(data signup.Data) {
	dataByte, _ := json.Marshal(&data)
	client := &http.Client{}
	url := "http://127.0.0.1:3004/v1/signup"
	method := "POST"

	req, err := http.NewRequest(method, url, bytes.NewReader(dataByte))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(body))
}

func formData(n int) signup.Data {
	var requestBody signup.Data
	requestBody.Login = fmt.Sprintf("%s-%d", func() []byte {
		name := gofakeit.Username()
		if len(name) >= 10 {
			return []byte(name)[0:10]
		} else {
			return []byte(name)
		}
	}(), n)
	requestBody.Password = gofakeit.Password(true, true, true, true, true, 20)
	requestBody.Age = gofakeit.UintRange(18, 100)
	requestBody.Hobbies = gofakeit.Sentence(100)
	requestBody.City = gofakeit.City()

	gender := gofakeit.Gender()
	switch gender {
	case "female":
		requestBody.Gender = "f"
		requestBody.Name = faker.FirstNameFemale()
		requestBody.Surname = faker.LastName()
	case "male":
		requestBody.Gender = "m"
		requestBody.Name = faker.FirstNameMale()
		requestBody.Surname = faker.LastName()
	}

	return requestBody
}
