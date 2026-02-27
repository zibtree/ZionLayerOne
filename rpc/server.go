package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zionlayer/zionlayer/core/mempool"
	"github.com/zionlayer/zionlayer/core/state"
	"github.com/zionlayer/zionlayer/core/transaction"
	"go.uber.org/zap"
)

// Request is a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

// Response is a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// RPCError represents a JSON-RPC error object.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Server is the ZionLayer JSON-RPC server.
type Server struct {
	state   *state.StateDB
	pool    *mempool.Pool
	logger  *zap.Logger
	port    int
}

// NewServer creates a new RPC server.
func NewServer(stateDB *state.StateDB, pool *mempool.Pool, logger *zap.Logger, port int) *Server {
	return &Server{state: stateDB, pool: pool, logger: logger, port: port}
}

// Start begins listening for RPC requests.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)
	mux.HandleFunc("/health", s.health)
	addr := fmt.Sprintf(":%d", s.port)
	s.logger.Info("RPC server starting", zap.String("addr", addr))
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, nil, -32700, "parse error")
		return
	}

	var result interface{}
	var rpcErr *RPCError

	switch req.Method {
	case "zion_getBalance":
		result, rpcErr = s.getBalance(req.Params)
	case "zion_sendTransaction":
		result, rpcErr = s.sendTransaction(req.Params)
	case "zion_getAgent":
		result, rpcErr = s.getAgent(req.Params)
	case "zion_getMempoolSize":
		result = map[string]int{"size": s.pool.Size()}
	case "zion_chainId":
		result = "0x1" // chain ID 1 for devnet
	default:
		rpcErr = &RPCError{Code: -32601, Message: "method not found"}
	}

	resp := Response{JSONRPC: "2.0", ID: req.ID, Result: result, Error: rpcErr}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) getBalance(params json.RawMessage) (interface{}, *RPCError) {
	var args []string
	if err := json.Unmarshal(params, &args); err != nil || len(args) == 0 {
		return nil, &RPCError{Code: -32602, Message: "invalid params"}
	}
	acc := s.state.GetAccount(args[0])
	return map[string]string{
		"address": acc.Address,
		"balance": acc.Balance.String(),
		"nonce":   fmt.Sprintf("%d", acc.Nonce),
	}, nil
}

func (s *Server) sendTransaction(params json.RawMessage) (interface{}, *RPCError) {
	var txs []*transaction.Tx
	if err := json.Unmarshal(params, &txs); err != nil || len(txs) == 0 {
		return nil, &RPCError{Code: -32602, Message: "invalid params"}
	}
	tx := txs[0]
	if err := s.pool.Add(tx); err != nil {
		return nil, &RPCError{Code: -32000, Message: err.Error()}
	}
	hash := tx.Hash()
	return fmt.Sprintf("0x%x", hash), nil
}

func (s *Server) getAgent(params json.RawMessage) (interface{}, *RPCError) {
	var args []string
	if err := json.Unmarshal(params, &args); err != nil || len(args) == 0 {
		return nil, &RPCError{Code: -32602, Message: "invalid params"}
	}
	rec, err := s.state.GetAgent(args[0])
	if err != nil {
		return nil, &RPCError{Code: -32000, Message: err.Error()}
	}
	return rec, nil
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func writeError(w http.ResponseWriter, id interface{}, code int, msg string) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &RPCError{Code: code, Message: msg},
	}
	json.NewEncoder(w).Encode(resp)
}
