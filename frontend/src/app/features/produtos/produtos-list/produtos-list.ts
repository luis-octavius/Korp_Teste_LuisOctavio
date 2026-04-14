import { Component, OnInit, inject } from '@angular/core';
import { Router } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';

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
  colunas = ['nome', 'saldo', 'acoes'];

  ngOnInit(): void {
    this.carregarProdutos();
  }

  carregarProdutos(): void {
    this.carregando = true;
    this.produtoService.listar().subscribe({
      next: (produtos) => {
        this.produtos = produtos;
        this.carregando = false;
      },
      error: () => {
        // Erro já tratado pelo interceptor global
        this.carregando = false;
      },
    });
  }

  novoProduto(): void {
    this.router.navigate(['/produtos/novo']);
  }
}
