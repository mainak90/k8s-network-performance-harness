package logging

import (
	"github.com/mainak90/colorify"
	"log"
	"os"
)

var colorStruct = colorify.Colorify{Attr: colorify.Bold, NoColor: "false"}

func InfoLog(in ...any) {
	log.New(os.Stdout, colorStruct.Sprintf(colorify.Blue, "INFO: "), log.Ldate|log.Ltime).Println(in...)
}

func WarnLog(in ...any) {
	log.New(os.Stdout, colorStruct.Sprintf(colorify.Yellow, "WARNING: "), log.Ldate|log.Ltime).Println(in...)
}

func ErrLog(in ...any) {
	log.New(os.Stdout, colorStruct.Sprintf(colorify.Red, "ERROR: "), log.Ldate|log.Ltime).Println(in...)
}
