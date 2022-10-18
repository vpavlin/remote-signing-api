package types

type INonceStorage interface {
	Store(n *NonceSerializable) error
	Load(chainId uint64, address string, contract *string) (*NonceSerializable, error)
}

type NonceSerializable struct {
	Address        string         `json:"address"`
	Contract       *string        `json:"contract,omitempty"`
	ChainId        uint64         `json:"chainId"`
	Nonce          uint64         `json:"nonce"`
	ReturnedNonces SortedNonceArr `json:"returnedNonces"`
	LastUsed       int64          `json:"lastUsed"`
}

type SortedNonceArr []uint64

func (arr SortedNonceArr) Less(i, j int) bool {
	return arr[i] < arr[j]
}

func (arr SortedNonceArr) Len() int { return len(arr) }

func (arr SortedNonceArr) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }

type IWalletStorage interface {
	Store(ws *WalletSerializable) error
	Load(address string, apiKeyHashed string) (*WalletSerializable, error)
}

type WalletSerializable struct {
	ApiKeyHashed string `json:"apiKeyHashed"`
	PublicKey    string `json:"publicKey"`
	PrivateKey   string `json:"privateKey"`
}
