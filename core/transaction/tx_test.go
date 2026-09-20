package transaction
import("crypto/ed25519";"crypto/rand";"math/big";"testing")
func TestTransactionSignature(t *testing.T){pub,priv,e:=ed25519.GenerateKey(rand.Reader);if e!=nil{t.Fatal(e)};tx:=NewTransferTx(AddressFromPublicKey(pub),"0xrecipient",big.NewInt(10),0,big.NewInt(1));if e=tx.Sign(priv);e!=nil{t.Fatal(e)};if e=tx.VerifySignature();e!=nil{t.Fatal(e)};tx.Value=big.NewInt(11);if e=tx.VerifySignature();e==nil{t.Fatal("tampered transaction verified")}}
