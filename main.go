package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/rudrodip/dummylsp/analysis"
	"github.com/rudrodip/dummylsp/lsp"
	"github.com/rudrodip/dummylsp/rpc"
)

func main() {
	logger := getLogger("dummylsp.log")
	logger.Println("Starting dummylsp")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error decoding message: %v", err)
			continue
		}
		handleMessage(logger, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, state analysis.State, method string, contents []byte) {
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
			logger.Printf("Error unmarshalling textDocument/didOpen request: %v", err)
			return
		}
		logger.Printf("Opened %s", request.Params.TextDocument.URI)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling textDocument/didChange request: %v", err)
			return
		}
		logger.Printf("Changed %s", request.Params.TextDocument.URI)

		for _, change := range request.Params.ContentChanges {
			state.OpenDocument(request.Params.TextDocument.URI, change.Text)
		}
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 06666)
	if err != nil {
		panic(err)
	}

	return log.New(logFile, "[dummylsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
