import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { catchError, throwError } from 'rxjs';

import {
  Nota,
  NotaDetalhe,
  AdicionarItensRequest,
  ImprimirNotaResponse,
} from '../models/nota.model';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class NotaService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = environment.api.faturamentoBaseUrl;

  criar(): Observable<Nota> {
    return this.http
      .post<Nota>(`${this.baseUrl}/notas`, {})
      .pipe(catchError((err) => throwError(() => err)));
  }

  listar(): Observable<Nota[]> {
    return this.http
      .get<Nota[]>(`${this.baseUrl}/notas`)
      .pipe(catchError((err) => throwError(() => err)));
  }

  buscarPorId(id: string): Observable<NotaDetalhe> {
    return this.http
      .get<NotaDetalhe>(`${this.baseUrl}/notas/${id}`)
      .pipe(catchError((err) => throwError(() => err)));
  }

  adicionarItens(notaId: string, req: AdicionarItensRequest): Observable<void> {
    return this.http
      .post<void>(`${this.baseUrl}/notas/${notaId}/itens`, req)
      .pipe(catchError((err) => throwError(() => err)));
  }

  imprimir(notaId: string): Observable<ImprimirNotaResponse> {
    return this.http
      .post<ImprimirNotaResponse>(`${this.baseUrl}/notas/${notaId}/imprimir`, {})
      .pipe(catchError((err) => throwError(() => err)));
  }
}
