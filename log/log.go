package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main()  {
	//panics()
	//news()
	//web()
	//http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {})
	//runtime.Breakpoint()

}
//go:embed manual.txt
func web()  {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(writer, request)
	})
	http.ListenAndServe(":801", nil)
	fmt.Println("-----")
}

// news 打印信息
func news ()  {
	infoLog := "info.log"
	create, err := os.Create(infoLog)
	defer create.Close()
	if err != nil {
		log.Fatalln("open file error", err)
	}
	logger := log.New(create, "[info]:", log.Llongfile)
	logger.Println("A Info message here")
	logger.SetPrefix("[debug]")
	logger.SetFlags(log.Lshortfile)
	logger.Println("A Debug message here")
}

// panics panic之前定义的defer才会执行
func panics ()  {
	defer func() {
		if err:= recover();err != nil {
			fmt.Println(err)
		}
	}()
	defer func() {
		fmt.Println("err")
	}()
	log.Panicln("report defer error")
}
