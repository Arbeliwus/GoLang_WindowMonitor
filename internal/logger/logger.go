package logger

import (
	"io"
	"log"
	"os"
)

var Error *log.Logger

func Init() {
	// 開啟/建立 error.log
	f, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("無法開啟 error.log:", err)
	}

	// 同時輸出到 console + 檔案
	multi := io.MultiWriter(os.Stdout, f)

	// 自訂錯誤 logger
	Error = log.New(multi, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
