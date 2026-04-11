package internal

type CriarNotaResponse struct {
	ID     string `json:"id"`
	NumSeq int64  `json:"num_seq"`
	Status string `json:"status"`
}

type ItemNotaRequest struct {
	ProdutoID  string `json:"produto_id"`
	Quantidade int32  `json:"quantidade"`
}

type AdicionarItemRequest struct {
	NotaID string            `json:"nota_id"`
	Items  []ItemNotaRequest `json:"items"`
}

type ItemNotaResponse struct {
	ID          string `json:"id"`
	ProdutoID   string `json:"produto_id"`
	Quantidade  int32  `json:"quantidade"`
	ProdutoNome string `json:"produto_nome"`
}

type NotaDetalheResponse struct {
	ID        string             `json:"id"`
	NumSeq    int64              `json:"num_seq"`
	Status    string             `json:"status"`
	CreatedAt string             `json:"created_at"`
	PrintedAt *string            `json:"printed_at,omitempty"`
	Items     []ItemNotaResponse `json:"items"`
}

type ImprimirNotaResponse struct {
	ID        string `json:"id"`
	NumSeq    int64  `json:"num_seq"`
	Status    string `json:"status"`
	PrintedAt string `json:"printed_at"`
}

// --------------- Estoque ------------------------------
type DebitarEstoqueRequest struct {
	ProdutoID  string `json:"produto_id"`
	Quantidade int32  `json:"quantidade"`
	NotaID     string `json:"nota_id"`
	NotaNum    int64  `json:"nota_num"`
}

type ReverterDebitoRequest struct {
	ProdutoID  string `json:"produto_id"`
	Quantidade int32  `json:"quantidade"`
	NotaID     string `json:"nota_id"`
	NotaNum    int64  `json:"nota_num"`
}

type EstoqueProdutoResponse struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Saldo int32  `json:"saldo"`
}
