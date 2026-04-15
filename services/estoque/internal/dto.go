package internal

type CriarProdutoRequest struct {
	Codigo string `json:"codigo"`
	Nome   string `json:"nome"`
	Saldo  int32  `json:"saldo"`
}

type ProdutoResponse struct {
	ID     string `json:"id"`
	Codigo string `json:"codigo"`
	Nome   string `json:"nome"`
	Saldo  int32  `json:"saldo"`
}

type DebitarEstoqueRequest struct {
	ProdutoID  string `json:"produto_id"`
	Quantidade int32  `json:"quantidade"`
	NotaID     string `json:"nota_id"`
	NotaNum    int64  `json:"nota_num"`
}

type DebitarEstoqueResponse struct {
	ProdutoID string `json:"produto_id"`
	NovoSaldo int32  `json:"novo_saldo"`
}

type ReverterDebitoRequest struct {
	ProdutoID  string `json:"produto_id"`
	Quantidade int32  `json:"quantidade"`
	NotaID     string `json:"nota_id"`
	NotaNum    int64  `json:"nota_num"`
}
