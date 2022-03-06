package calendar

import (
	"fmt"
	"os"
	"strings"
)

const data_path = "data/"

// interface

type Loader interface {
	Load(name string) (User, error)
	Save(user *User) error
}

// types
type PlainLoader struct {
}

func (loader *PlainLoader) Load(name string) (User, error) {
	convertToUser := func(data []byte) (User, error) {
		to_convert := string(data)

		user_data := strings.Split(to_convert, "\r\n")

		if len(user_data) <= 1 {
			return User{}, fmt.Errorf("could not convert %s", to_convert)
		}

		return User{Name: name, Mail: user_data[0], Dates: user_data[1:]}, nil
	}
	content, err := os.ReadFile(data_path + name + ".txt")

	if err != nil {
		return User{Name: name, Mail: ""}, err
	}
	return convertToUser(content)
}

func (loader *PlainLoader) Save(user *User) error {
	filename := data_path + user.Name + ".txt"
	data := []byte(user.Mail + "\r\n")
	data = append(data, []byte(strings.Join(user.Dates, "\r\n"))...)
	return os.WriteFile(filename, data, 0600)
}

func LoaderFactory(loader string) Loader {
	if loader == "file" {
		return &PlainLoader{}
	}

	return &PlainLoader{}
}

func CreateUser(name string, mail string, dates string) User {
	dates_array := strings.Split(dates, "\r\n")
	return User{Name: name, Mail: mail, Dates: dates_array}
}

type User struct {
	Name  string
	Mail  string
	Dates []string
}
