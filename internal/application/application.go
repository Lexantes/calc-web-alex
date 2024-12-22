package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"calc-web-alex/pkg/calc"
)

type Config struct {
	Addr string
}

func ConfigFromLine() *Config {
	config := new(Config)
	for i := 0; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--port":
			config.Addr = os.Args[i+1]
		}
	}
	if config.Addr == "" {
		config.Addr = "30001"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromLine(),
	}
}

func (a *Application) Run() error {
	for {
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}
		if text == "exit" {
			log.Println("aplication was successfully closed")
			return nil
		}
		result, err := calc.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := calc.Calc(request.Expression)
	if err != nil {
		fmt.Fprintf(w, "err: %s", err.Error())
	} else {
		fmt.Fprintf(w, "result: %f", result)
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
