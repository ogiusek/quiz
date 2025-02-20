package usersmodule

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

// user name
// user password

// user name

type UserName string

func (u *UserName) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)

	if err != nil {
		return err
	}

	*u = UserName(strings.Trim(strings.ToLower(s), " "))

	return nil
}

func (UserName) GormDataType() string { return "varchar(64)" }

var nameRegex = regexp.MustCompile("^[a-zA-Z0-9-_]*$")

func (name *UserName) Valid() []error {
	var errs []error
	if len(*name) < 3 {
		errs = append(errs, errors.New("name has to be at least 3 characters long"))
	}
	if len(*name) > 64 {
		errs = append(errs, errors.New("name cannot have more than 64 characters"))
	}
	if ok := nameRegex.MatchString(string(*name)); !ok {
		errs = append(errs, errors.New("name can contain only letters, numbers, - and _"))
	}
	return errs
}

// user password

type UserPassword string

var smallLeterRegex = regexp.MustCompile("[a-z]")

// var bigLeterRegex = regexp.MustCompile("[A-Z]")
var numberRegex = regexp.MustCompile("[0-9]")

// var specialCharacterRegex = regexp.MustCompile("[!@#$%^&*()_+{}:\"|<>?]")

func (password *UserPassword) Valid() []error {
	var errs []error
	if len(*password) < 3 {
		errs = append(errs, errors.New("password has to be at least 3 characters long"))
	}
	if len(*password) > 64 {
		errs = append(errs, errors.New("password cannot have more than 64 characters"))
	}
	if ok := smallLeterRegex.MatchString(string(*password)); !ok {
		errs = append(errs, errors.New("password must contain small letter"))
	}
	// if ok := bigLeterRegex.MatchString(string(*password)); !ok {
	// 	errs = append(errs, errors.New("password must contain big letter"))
	// }
	if ok := numberRegex.MatchString(string(*password)); !ok {
		errs = append(errs, errors.New("password must contain number"))
	}
	// if ok := specialCharacterRegex.MatchString(string(*password)); !ok {
	// 	errs = append(errs, errors.New("password must contain special character"))
	// }
	return errs
}

func (UserPassword) GetValue() (driver.Value, error) {
	return nil, errors.New("cannot store raw password in database")
	// return string(password), nil
}
