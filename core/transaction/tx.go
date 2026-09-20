// Package transaction defines signed, chain-bound ZionLayer transactions.\n// Protocol baseline: September 2026.\npackage transaction\n
import ("crypto/ed25519";"crypto/sha256";"encoding/hex";"encoding/json";"errors";"math/big")
type TxType uint8
const (TxTransfer TxType=iota;TxAgentRegister;TxAgentMessage;TxAgentDelegate;TxDeployContract;TxCallContract;TxInferenceReceipt;TxValidatorStake;TxValidatorUnstake;TxA2HPost;TxA2HClaim;TxA2HComplete)
const ChainID uint64=1
var(ErrInvalidSignature=errors.New("invalid transaction signature");ErrInvalidPublicKey=errors.New("invalid public key");ErrInvalidAddress=errors.New("public key does not match sender address");ErrInvalidValue=errors.New("transaction value must be non-negative"))
type Capability struct{Name string;Version string}
type AgentDID struct{ID string;Controller string;Capabilities []Capability;PublicKey []byte;Metadata map[string]string}
type MessageType string
const(MsgTask MessageType="TASK";MsgResult MessageType="RESULT";MsgDelegate MessageType="DELEGATE";MsgRevoke MessageType="REVOKE")
type AgentMessage struct{From string;To string;Type MessageType;Payload []byte;Nonce uint64}
type InferenceReceipt struct{AgentID string;ModelHash []byte;InputHash []byte;OutputHash []byte;Timestamp int64;ProverSig []byte}
type A2HTaskStatus uint8
const(A2HOpen A2HTaskStatus=iota;A2HClaimed;A2HComplete;A2HDisputed;A2HExpired)
type A2HTask struct{ID string;AgentID string;Title string;Description string;Skills []string;Reward *big.Int;Deadline uint64;Assignee string;Status A2HTaskStatus;CreatedAt uint64}
func(t *A2HTask)ComputeID()string{b,_:=json.Marshal(struct{AgentID,Title,Description string;Skills []string;Reward *big.Int;Deadline uint64}{t.AgentID,t.Title,t.Description,t.Skills,t.Reward,t.Deadline});h:=sha256.Sum256(b);return "task_"+hex.EncodeToString(h[:])}
type Tx struct{Type TxType;From string;To string;Value *big.Int;Gas uint64;GasPrice *big.Int;Nonce uint64;ChainID uint64;Data json.RawMessage;PublicKey []byte;Signature []byte}
func(tx *Tx)Hash()[32]byte{cp:=*tx;cp.Signature=nil;b,_:=json.Marshal(cp);return sha256.Sum256(b)}
func AddressFromPublicKey(pub ed25519.PublicKey)string{h:=sha256.Sum256(pub);return "0x"+hex.EncodeToString(h[len(h)-20:])}
func(tx *Tx)SigningBytes()[]byte{h:=tx.Hash();return h[:]}
func(tx *Tx)Sign(priv ed25519.PrivateKey)error{if len(priv)!=ed25519.PrivateKeySize{return ErrInvalidPublicKey};pub:=priv.Public().(ed25519.PublicKey);tx.PublicKey=append([]byte(nil),pub...);if tx.From==""{tx.From=AddressFromPublicKey(pub)};if tx.From!=AddressFromPublicKey(pub){return ErrInvalidAddress};tx.Signature=ed25519.Sign(priv,tx.SigningBytes());return nil}
func(tx *Tx)VerifySignature()error{if len(tx.PublicKey)!=ed25519.PublicKeySize||len(tx.Signature)!=ed25519.SignatureSize{return ErrInvalidSignature};pub:=ed25519.PublicKey(tx.PublicKey);if tx.From!=AddressFromPublicKey(pub){return ErrInvalidAddress};if !ed25519.Verify(pub,tx.SigningBytes(),tx.Signature){return ErrInvalidSignature};if tx.ChainID!=ChainID{return errors.New("wrong chain id")};if tx.Value!=nil&&tx.Value.Sign()<0{return ErrInvalidValue};return nil}
func NewTransferTx(from,to string,value *big.Int,nonce uint64,gasPrice *big.Int)*Tx{return &Tx{Type:TxTransfer,From:from,To:to,Value:new(big.Int).Set(value),Gas:21000,GasPrice:new(big.Int).Set(gasPrice),Nonce:nonce,ChainID:ChainID}}
func NewAgentRegisterTx(from string,did AgentDID,nonce uint64,gasPrice *big.Int)*Tx{b,_:=json.Marshal(did);return &Tx{Type:TxAgentRegister,From:from,Gas:200000,GasPrice:new(big.Int).Set(gasPrice),Nonce:nonce,ChainID:ChainID,Data:b}}
func NewAgentMessageTx(from string,msg AgentMessage,nonce uint64,gasPrice *big.Int)*Tx{b,_:=json.Marshal(msg);return &Tx{Type:TxAgentMessage,From:from,Gas:50000,GasPrice:new(big.Int).Set(gasPrice),Nonce:nonce,ChainID:ChainID,Data:b}}
func NewInferenceReceiptTx(from string,r InferenceReceipt,nonce uint64,gasPrice *big.Int)*Tx{b,_:=json.Marshal(r);return &Tx{Type:TxInferenceReceipt,From:from,Gas:100000,GasPrice:new(big.Int).Set(gasPrice),Nonce:nonce,ChainID:ChainID,Data:b}}
