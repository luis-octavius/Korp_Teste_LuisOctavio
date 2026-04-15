export interface Produto {
  id: string;
  codigo: string;
  nome: string;
  saldo: number;
}

export interface CriarProdutoRequest {
  codigo: string;
  nome: string;
  saldo: number;
}
