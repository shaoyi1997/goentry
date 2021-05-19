package logger

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	isProduction  = flag.Bool("p", false, "Specifies the environment of the application")
	filename      = filepath.Base(os.Args[0])
)

func InitLogger() {
	flag.Parse()

	out := getOutputWriter()

	InfoLogger = log.New(out, makePrefix("INFO"), log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(out, makePrefix("WARNING"), log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(out, makePrefix("ERROR"), log.Ldate|log.Ltime|log.Lshortfile)
}

func makePrefix(prefix string) string {
	return fmt.Sprintf("[%s] %s: ", filename, prefix)
}

func getOutputWriter() io.Writer {
	if !*isProduction {
		return os.Stdout
	}

	file, err := os.OpenFile(filename+"-logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return io.Writer(file)
}
