import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { catchError, throwError } from 'rxjs';

import { Produto, CriarProdutoRequest } from '../models/produto.model';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class ProdutoService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = environment.api.estoqueBaseUrl;

  listar(): Observable<Produto[]> {
    return this.http
      .get<Produto[]>(`${this.baseUrl}/produtos`)
      .pipe(catchError((err) => throwError(() => err)));
  }

  buscarPorId(id: string): Observable<Produto> {
    return this.http
      .get<Produto>(`${this.baseUrl}/produtos/${id}`)
      .pipe(catchError((err) => throwError(() => err)));
  }

  criar(req: CriarProdutoRequest): Observable<Produto> {
    return this.http
      .post<Produto>(`${this.baseUrl}/produtos`, req)
      .pipe(catchError((err) => throwError(() => err)));
  }
}
