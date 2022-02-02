package embedded

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
)

//go:embed *.rego
var EmbeddedFiles embed.FS

func GetFiles() {

	bt, err := fs.ReadFile(EmbeddedFiles, "Allow-Image-Name.rego")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bt))
}
