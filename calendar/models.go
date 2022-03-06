package calendar

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const data_path = "data/"

// interface

type Loader interface {
	Load(name string) (User, error)
	Save(user *User) error
}

// types
type PlainLoader struct {
	Users map[string]User
}

func (loader *PlainLoader) createUser(name string, data []byte) error {
	to_convert := string(data)

	user_data := strings.Split(to_convert, "\r\n")

	if len(user_data) <= 1 {
		return fmt.Errorf("could not convert %s", to_convert)
	} else if len(user_data) >= 3 {
		loader.Users[name] = User{
			Name:     name,
			Mail:     user_data[0],
			Password: []byte(user_data[1]),
			Salt:     []byte(user_data[2]),
			Dates:    user_data[3:],
		}
	}
	return nil
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
	Name     string
	Mail     string
	Password []byte
	Salt     []byte
	Dates    []string
}

type UserSession struct {
	UserId          int32
	SessionKey      []byte
	LoginTime       time.Time
	LastInteraction time.Time
}

type UserSessions struct {
	activeUsers map[string]UserSession
}

func CreateSession() *UserSessions {
	return &UserSessions{activeUsers: make(map[string]UserSession)}
}
