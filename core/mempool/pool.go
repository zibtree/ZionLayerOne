package mempool
import("errors";"sort";"sync";"github.com/zionlayer/zionlayer/core/transaction")
const MaxPoolSize=10000
var(ErrPoolFull=errors.New("mempool is full");ErrDuplicateTx=errors.New("duplicate transaction"))
type Pool struct{mu sync.RWMutex;txs map[[32]byte]*transaction.Tx}
func NewPool()*Pool{return &Pool{txs:map[[32]byte]*transaction.Tx{}}}
func(p *Pool)Add(tx *transaction.Tx)error{if tx==nil{return errors.New("nil transaction")};if tx.Gas==0{return errors.New("gas limit required")};if tx.GasPrice==nil||tx.GasPrice.Sign()<0{return errors.New("invalid gas price")};if e:=tx.VerifySignature();e!=nil{return e};p.mu.Lock();defer p.mu.Unlock();if len(p.txs)>=MaxPoolSize{return ErrPoolFull};h:=tx.Hash();if _,ok:=p.txs[h];ok{return ErrDuplicateTx};p.txs[h]=tx;return nil}
func(p *Pool)Pop(n int)[]*transaction.Tx{p.mu.Lock();defer p.mu.Unlock();a:=make([]*transaction.Tx,0,len(p.txs));for _,t:=range p.txs{a=append(a,t)};sort.Slice(a,func(i,j int)bool{return a[i].GasPrice.Cmp(a[j].GasPrice)>0});if n<0{n=0};if n>len(a){n=len(a)};out:=a[:n];for _,t:=range out{delete(p.txs,t.Hash())};return out}
func(p *Pool)Size()int{p.mu.RLock();defer p.mu.RUnlock();return len(p.txs)}
