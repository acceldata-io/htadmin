// Acceldata Inc. and its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// 	Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package htpasswd is utility package to manipulate htpasswd files. I supports\
// bcrypt and sha hashes.
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

// HashedPasswords name => hash
type HashedPasswords map[string]string

// HashAlgorithm enum for hashing algorithms
type HashAlgorithm string

type htdata struct {
	File    string  `json:"file"`
	Account account `json:"account"`
}

const (
	// HashAPR1 Apache MD5 crypt - legacy
	HashAPR1 = "apr1"
	// HashBCrypt bcrypt - recommended
	HashBCrypt = "bcrypt"
	// HashSHA sha5 insecure - do not use
	HashSHA = "sha"
)

const (
	// PasswordSeparator separates passwords from hashes
	PasswordSeparator = ":"
	// LineSeparator separates password records
	LineSeparator = "\n"
)

// MaxHtpasswdFilesize if your htpassd file is larger than 8MB, then your are doing it wrong
const MaxHtpasswdFilesize = 8 * 1024 * 1024

// ErrNotExist is the error returned when a user does not exist.
var ErrNotExist = errors.New("user did not exist in file")

// Bytes bytes representation
func (hp HashedPasswords) Bytes() (passwordBytes []byte) {
	passwordBytes = []byte{}
	for name, hash := range hp {
		passwordBytes = append(passwordBytes, []byte(name+PasswordSeparator+hash+LineSeparator)...)
	}
	return passwordBytes
}

// WriteToFile put them to a file will be overwritten or created
func (hp HashedPasswords) WriteToFile(file string) error {
	return ioutil.WriteFile(file, hp.Bytes(), 0o644)
}

// SetPassword set a password for a user with a hashing algo
func (hp HashedPasswords) SetPassword(name, password string, hashAlgorithm HashAlgorithm) (err error) {
	if len(password) == 0 {
		return errors.New("passwords must not be empty, if you want to delete a user call deleteUser")
	}

	hash := ""
	prefix := ""
	switch hashAlgorithm {
	case HashAPR1:
		hash, err = hashApr1(password)
	case HashBCrypt:
		hash, err = hashBcrypt(password)
	case HashSHA:
		prefix = "{SHA}"
		hash = hashSha(password)
	}
	if err != nil {
		return err
	}
	hp[name] = prefix + hash
	return nil
}

// DeleteUser deletes the user from the file
func (hp HashedPasswords) DeleteUser(name string) (err error) {
	delete(hp, name)
	return nil
}

// Check if the user exists
func (hp HashedPasswords) isUserExists(name string) bool {
	//
	if _, ok := hp[name]; ok {
		// User already exists
		return true
	}
	return false
}

// parseHtpasswdFile load a htpasswd file
func parseHtpasswdFile(file string) (passwords HashedPasswords, err error) {
	htpasswdBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	if len(htpasswdBytes) > MaxHtpasswdFilesize {
		err = errors.New("this file is too large, use a database instead")
		return
	}
	return parseHtpasswd(htpasswdBytes)
}

// parseHtpasswd parse htpasswd bytes
func parseHtpasswd(htpasswdBytes []byte) (passwords HashedPasswords, err error) {
	lines := strings.Split(string(htpasswdBytes), LineSeparator)
	passwords = make(map[string]string)
	for lineNumber, line := range lines {
		// scan lines
		line = strings.Trim(line, " ")
		if len(line) == 0 {
			// skipping empty lines
			continue
		}
		parts := strings.Split(line, PasswordSeparator)
		if len(parts) != 2 {
			err = errors.New(fmt.Sprintln("invalid line", lineNumber+1, "unexpected number of parts split by", PasswordSeparator, len(parts), "instead of 2 in\"", line, "\""))
			return
		}
		for i, part := range parts {
			parts[i] = strings.Trim(part, " ")
		}
		_, alreadyExists := passwords[parts[0]]
		if alreadyExists {
			err = errors.New("invalid htpasswords file - user " + parts[0] + " was already defined")
			return
		}
		passwords[parts[0]] = parts[1]
	}
	return
}

// SetPassword set password for a user with a given hashing algorithm
func (ht htdata) SetPassword() error {
	_, err := os.Stat(ht.File)
	passwords := HashedPasswords(map[string]string{})
	if err == nil {
		passwords, err = parseHtpasswdFile(ht.File)
		if err != nil {
			return err
		}
	}
	if !passwords.isUserExists(ht.Account.Name) {
		err = passwords.SetPassword(ht.Account.Name, ht.Account.Password, HashAPR1)
		if err != nil {
			return err
		}
		//
		return passwords.WriteToFile(ht.File)
	}

	return errors.New("user already exists")
}

// DeleteUser deletes the user from the file
func (ht htdata) DeleteUser() error {
	_, err := os.Stat(ht.File)
	passwords := HashedPasswords(map[string]string{})
	if err == nil {
		passwords, err = parseHtpasswdFile(ht.File)
		if err != nil {
			return err
		}
	}
	if passwords.isUserExists(ht.Account.Name) {
		err = passwords.DeleteUser(ht.Account.Name)
		if err != nil {
			return err
		}
		//
		return passwords.WriteToFile(ht.File)
	}

	return errors.New("user doesn't exists")
}

func randSeq(n int) string {
	letters := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
