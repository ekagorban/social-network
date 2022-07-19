package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	"github.com/bxcodec/faker/v3"
)

var genders = []string{"f", "m"}

func main() {
	url := "http://127.0.0.1:3004/v1/signup"
	method := "POST"
	gender := genders[rand.Int()%len(genders)]
	switch gender {
	case "f":

		password := faker.Password()
		name := faker.FirstNameFemale()
		surname := faker.LastName()
		//login:=fmt.Sprintf("%s-%s-%s",surname,name,time.Now().Unix())
		login := faker.Username()
	//	age:=faker.RandomInt() rand.Intn()

	case "m":
	}

	payload := strings.NewReader(`{
    "login":"login_user1",
    "password":"password1",
    "name":"name_user1",
    "surname":"surname_user1",
    "age":20,
    "gender":"m",
    "hobbies":"hobbies1",
    "city":"city1"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
