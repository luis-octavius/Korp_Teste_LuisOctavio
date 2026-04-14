import { Component, OnInit, inject } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatTableModule } from '@angular/material/table';
import { MatDividerModule } from '@angular/material/divider';
import { MatDialogModule, MatDialog } from '@angular/material/dialog';
import { DatePipe } from '@angular/common';

import { NotaService } from '../../../core/services/nota.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner';
import { NotaDetalhe } from '../../../core/models/nota.model';
import { finalize } from 'rxjs';

@Component({
  selector: 'app-nota-detalhe',
  standalone: true,
  imports: [
    DatePipe,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatChipsModule,
    MatTableModule,
    MatDividerModule,
    MatDialogModule,
    LoadingSpinnerComponent,
  ],
  templateUrl: './nota-detalhe.html',
  styleUrl: './nota-detalhe.scss',
})
export class NotaDetalheComponent implements OnInit {
  private readonly notaService = inject(NotaService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly dialog = inject(MatDialog);

  nota: NotaDetalhe | null = null;
  carregando = false;
  imprimindo = false;
  colunas = ['produto_nome', 'quantidade'];

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id') ?? '';
    this.carregarNota(id);
  }

  carregarNota(id: string): void {
    this.carregando = true;
    this.notaService.buscarPorId(id).subscribe({
      next: (nota) => {
        this.nota = nota;
        this.carregando = false;
      },
      error: () => {
        this.carregando = false;
      },
    });
  }

  adicionarItens(): void {
    this.router.navigate(['/notas', this.nota!.id, 'itens']);
  }

  imprimir(): void {
    if (!this.nota || this.nota.status !== 'ABERTA') return;

    this.imprimindo = true;

    this.notaService
      .imprimir(this.nota.id)
      .pipe(
        // finalize sempre executa — seja sucesso ou erro
        finalize(() => (this.imprimindo = false)),
      )
      .subscribe({
        next: (resultado) => {
          // Atualiza o status localmente sem precisar recarregar
          this.nota = {
            ...this.nota!,
            status: resultado.status,
            printed_at: resultado.printed_at,
          };
        },
        error: () => {
          // Erro já tratado pelo interceptor global via snackbar
        },
      });
  }

  voltar(): void {
    this.router.navigate(['/notas']);
  }

  get totalItens(): number {
    return this.nota?.items.reduce((acc, item) => acc + item.quantidade, 0) ?? 0;
  }
}
