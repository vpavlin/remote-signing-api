package types

type INonceStorage interface {
	Store(n *NonceSerializable) error
	Load(chainId uint64, address string) (*NonceSerializable, error)
}

type NonceSerializable struct {
	Address        string         `json:"address"`
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
