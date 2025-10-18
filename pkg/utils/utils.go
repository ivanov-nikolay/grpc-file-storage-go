package utils

import "os"

func GetProjectRoot() (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return root, nil
}
