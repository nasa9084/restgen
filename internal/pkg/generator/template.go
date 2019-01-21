package generator

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nasa9084/restgen/internal/pkg/assets"
)

func Template(src, dest string, values ...interface{}) (string, error) {
	srcFile, err := assets.Assets.Open(src)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	tmpl, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return "", err
	}

	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	out := fmt.Sprintf(string(tmpl), values...)

	if _, err := fmt.Fprint(destFile, out); err != nil {
		return "", err
	}
	return out, nil
}
