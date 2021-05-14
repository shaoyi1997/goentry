package logger

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	isProduction = flag.Bool("P", false, "Specifies the environment of the application")
)

func init() {
	flag.Parse()

	out := getOutputWriter()

	InfoLogger = log.New(out, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(out, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func getOutputWriter() io.Writer {
	if !*isProduction {
		return os.Stdout
	}

	filename := filepath.Base(os.Args[0])
	file, err := os.OpenFile(filename + "-logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return io.Writer(file)
}