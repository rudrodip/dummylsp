package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/rudrodip/dummylsp/lsp"
	"github.com/rudrodip/dummylsp/rpc"
)

func main() {
	logger := getLogger("dummylsp.log")
	logger.Println("Starting dummylsp")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error decoding message: %v", err)
			continue
		}
		handleMessage(logger, method, contents)
	}
}

func handleMessage(logger *log.Logger, method string, contents []byte) {
	logger.Printf("Receive message with method %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling initialize request: %v", err)
		}
		logger.Printf("Client info: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)

		writer := os.Stdout
		writer.Write([]byte(reply))
		logger.Print("Sent response")

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling initialize request: %v", err)
		}
		logger.Printf("Opened %s %s", request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 06666)
	if err != nil {
		panic(err)
	}

	return log.New(logFile, "[dummylsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
