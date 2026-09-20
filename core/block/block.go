package block
import("crypto/sha256";"encoding/json";"time";"github.com/zionlayer/zionlayer/core/transaction")
type Header struct{Version uint32;Height uint64;Timestamp int64;PrevHash [32]byte;StateRoot [32]byte;TxRoot [32]byte;AgentRoot [32]byte;ValidatorAddr []byte;Signature []byte}
type Block struct{Header Header;Txs []*transaction.Tx}
func NewBlock(h uint64,p [32]byte,v []byte,t []*transaction.Tx)*Block{return &Block{Header:Header{Version:1,Height:h,Timestamp:time.Now().UnixNano(),PrevHash:p,ValidatorAddr:append([]byte(nil),v...)},Txs:t}}
func TxRoot(txs []*transaction.Tx)[32]byte{if len(txs)==0{return sha256.Sum256(nil)};l:=make([][32]byte,len(txs));for i,t:=range txs{l[i]=t.Hash()};for len(l)>1{n:=make([][32]byte,0,(len(l)+1)/2);for i:=0;i<len(l);i+=2{r:=l[i];if i+1<len(l){r=l[i+1]};var b [64]byte;copy(b[:32],l[i][:]);copy(b[32:],r[:]);n=append(n,sha256.Sum256(b[:]))};l=n};return l[0]}
func(b *Block)Finalize(root [32]byte){b.Header.StateRoot=root;b.Header.TxRoot=TxRoot(b.Txs);b.Header.AgentRoot=root}
func(b *Block)Hash()[32]byte{b2:=b.Header;b2.Signature=nil;d,_:=json.Marshal(b2);return sha256.Sum256(d)}
func GenesisBlock()*Block{return &Block{Header:Header{Version:1,Height:0,Timestamp:time.Date(2025,1,1,0,0,0,0,time.UTC).UnixNano(),TxRoot:sha256.Sum256(nil),StateRoot:sha256.Sum256(nil),AgentRoot:sha256.Sum256(nil)},Txs:[]*transaction.Tx{}}}
