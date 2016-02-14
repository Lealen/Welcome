package main

import "io/ioutil"

func saveToFile(file, what string) error {
	return ioutil.WriteFile(file, []byte(what), 0644)
}

func loadFromFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return string(data), err
}
