package main

import (
	"bufio"
	"encoding/json"
	"io"
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

	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error decoding message: %v", err)
			continue
		}
		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
	logger.Printf("Receive message with method %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling initialize request: %v", err)
		}
		logger.Printf("Client info: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, msg)
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

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling textDocument/hover request: %v", err)
			return
		}

		// create a response
		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		// write back
		writeResponse(writer, response)
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 06666)
	if err != nil {
		panic(err)
	}

	return log.New(logFile, "[dummylsp]", log.Ldate|log.Ltime|log.Lshortfile)
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}
