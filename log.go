package tcpx

import (
	"io"
	"log"
	"os"
)

var (
	//debug信息
	Debug *log.Logger = log.New(os.Stdout, "\u001B[1;36m[Debug]:\u001B[0m", log.Ltime|log.Llongfile)
	//重要信息
	Info *log.Logger = log.New(os.Stdout, "\u001B[1;34m[Info]:\u001B[0m", log.Ldate|log.Ltime|log.Lshortfile)
	//警告
	Warn *log.Logger = log.New(os.Stdout, "\u001B[1;33m[Warn]:\u001B[0m", log.Ldate|log.Ltime|log.Lshortfile)
	//错误
	Error *log.Logger
)

func ErrorInit(fileName string) {
	file, err := os.OpenFile(fileName+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("无法打开错误的log文件: ", err)
	}
	Error = log.New(io.MultiWriter(file, os.Stderr), "\u001B[1;31m[Error]:\u001B[0m", log.Ldate|log.Ltime|log.Lshortfile)
}

func tcpxPrintln(s string) {
	println("\x1b[1;34m" + s + "\x1b[0m")
}

func Fatal(s string) {
	log.Fatalln(s)
}
