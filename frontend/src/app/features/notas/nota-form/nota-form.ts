import { Component, OnInit, inject } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormBuilder, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatSelectModule } from '@angular/material/select';

import { NotaService } from '../../../core/services/nota.service';
import { ProdutoService } from '../../../core/services/produto.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner';
import { Produto } from '../../../core/models/produto.model';
import { ItemNotaRequest } from '../../../core/models/nota.model';

@Component({
  selector: 'app-nota-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatSelectModule,
    LoadingSpinnerComponent,
  ],
  templateUrl: './nota-form.html',
  styleUrl: './nota-form.scss',
})
export class NotaFormComponent implements OnInit {
  private readonly notaService = inject(NotaService);
  private readonly produtoService = inject(ProdutoService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly fb = inject(FormBuilder);

  notaId = '';
  produtos: Produto[] = [];
  itens: ItemNotaRequest[] = [];
  salvando = false;
  carregandoProdutos = false;

  itemForm = this.fb.group({
    produto_id: ['', Validators.required],
    quantidade: [1, [Validators.required, Validators.min(1)]],
  });

  ngOnInit(): void {
    this.notaId = this.route.snapshot.paramMap.get('id') ?? '';
    this.carregarProdutos();
  }

  carregarProdutos(): void {
    this.carregandoProdutos = true;
    this.produtoService.listar().subscribe({
      next: (produtos) => {
        this.produtos = produtos;
        this.carregandoProdutos = false;
      },
      error: () => {
        this.carregandoProdutos = false;
      },
    });
  }

  adicionarItem(): void {
    if (this.itemForm.invalid) return;

    const { produto_id, quantidade } = this.itemForm.value;

    // Evita duplicar o mesmo produto
    const jaAdicionado = this.itens.some((i) => i.produto_id === produto_id);
    if (jaAdicionado) {
      this.itemForm.get('produto_id')?.setErrors({ duplicado: true });
      return;
    }

    this.itens.push({ produto_id: produto_id!, quantidade: quantidade! });
    this.itemForm.reset({ quantidade: 1 });
  }

  removerItem(index: number): void {
    this.itens.splice(index, 1);
  }

  nomeProduto(id: string): string {
    return this.produtos.find((p) => p.id === id)?.nome ?? id;
  }

  salvar(): void {
    if (this.itens.length === 0) return;

    this.salvando = true;
    this.notaService.adicionarItens(this.notaId, { items: this.itens }).subscribe({
      next: () => {
        this.router.navigate(['/notas', this.notaId]);
      },
      error: () => {
        this.salvando = false;
      },
    });
  }

  cancelar(): void {
    this.router.navigate(['/notas', this.notaId]);
  }
}
