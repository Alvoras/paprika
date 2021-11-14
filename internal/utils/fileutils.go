package utils

import (
	"io/ioutil"
)

func CopyFile(src string, dst string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
