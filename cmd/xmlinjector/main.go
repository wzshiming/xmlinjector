package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/wzshiming/xmlinjector"
)

var (
	file        = "-"
	output      = "-"
	injectorKey = "INJECTOR"
	dataKey     = "data"
)

func init() {
	flag.StringVar(&file, "f", file, "read from file")
	flag.StringVar(&output, "o", output, "output to file")
	flag.StringVar(&injectorKey, "k", injectorKey, "key of inject")
	flag.StringVar(&dataKey, "d", dataKey, "data of args")
	flag.Parse()
}

func main() {
	data, err := open(file)
	if err != nil {
		log.Fatal(err)
	}

	out, err := xmlinjector.Inject([]byte(injectorKey), data, func(args, origin []byte) []byte {
		arg := xmlinjector.NewArgs(string(args), true)
		data, _ := arg.Get(dataKey)
		return []byte(data)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = write(output, out)
	if err != nil {
		log.Fatal(err)
	}
}

func open(filename string) ([]byte, error) {
	if filename == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(filename)
}

func write(filename string, data []byte) error {
	if filename == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(filename, data, 0544)
}
