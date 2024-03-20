package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

type ServerCapabilities struct {
	TextDocumentSync   int  `json:"textDocumentSync"`
	HoverProvider      bool `json:"hoverProvider,omitempty"`
	DefinitionProvider bool `json:"definitionProvider,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
			},
			ServerInfo: &ServerInfo{
				Name:    "dummylsp",
				Version: "0.0.1-alpha",
			},
		},
	}
}
