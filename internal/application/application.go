package application

import (
	"encoding/json"
	"fmt"
	"io"
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

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request struct {
		Expression string `json:"expression"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil || request.Expression == "" {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusBadRequest)
		return
	}

	result, err := calc.Calc(request.Expression)
	if err != nil {
		var errorMsg string
		statusCode := http.StatusUnprocessableEntity

		switch err {
		case calc.ErrInvalidExpression:
			errorMsg = "Error calculation"
		case calc.ErrDivisionByZero:
			errorMsg = "Division by zero"
		case calc.ErrMismatchedParentheses:
			errorMsg = "Mismatched parentheses"
		case calc.ErrUnexpectedToken:
			errorMsg = "Unexpected token"
		case calc.ErrNotEnoughValues:
			errorMsg = "Not enough values"
		case calc.ErrEmptyInput:
			errorMsg = "Empty input"
		default:
			errorMsg = "Error calculation"
			statusCode = http.StatusUnprocessableEntity
		}

		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, errorMsg), statusCode)
		return
	}

	response := struct {
		Result string `json:"result"`
	}{
		Result: fmt.Sprintf("%v", result),
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error while marshaling response: %v", err)
		http.Error(w, `{"error":"Unknown error occurred"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJson)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"Not Found"}`, http.StatusNotFound)
	})
	log.Printf("Server started on %s port", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
