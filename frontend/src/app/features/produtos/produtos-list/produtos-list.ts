import { Component, OnInit, inject } from '@angular/core';
import { Router } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { finalize, take } from 'rxjs';

import { ProdutoService } from '../../../core/services/produto.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner';
import { Produto } from '../../../core/models/produto.model';

@Component({
  selector: 'app-produtos-list',
  standalone: true,
  imports: [MatTableModule, MatButtonModule, MatIconModule, MatCardModule, LoadingSpinnerComponent],
  templateUrl: './produtos-list.html',
  styleUrl: './produtos-list.scss',
})
export class ProdutosListComponent implements OnInit {
  private readonly produtoService = inject(ProdutoService);
  private readonly router = inject(Router);

  produtos: Produto[] = [];
  carregando = false;
  colunas = ['codigo', 'nome', 'saldo'];

  ngOnInit(): void {
    this.carregarProdutos();
  }

  carregarProdutos(): void {
    this.carregando = true;
    this.produtoService
      .listar()
      .pipe(
        take(1),
        finalize(() => {
          this.carregando = false;
        }),
      )
      .subscribe({
        next: (produtos) => {
          // Always assign a new reference so table updates reliably.
          this.produtos = Array.isArray(produtos) ? [...produtos] : [];
        },
      });
  }

  novoProduto(): void {
    this.router.navigate(['/produtos/novo']);
  }
}
