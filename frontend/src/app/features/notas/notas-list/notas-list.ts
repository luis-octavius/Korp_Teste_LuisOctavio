import { Component, OnInit, inject } from '@angular/core';
import { Router } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';

import { NotaService } from '../../../core/services/nota.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner';
import { Nota } from '../../../core/models/nota.model';

@Component({
  selector: 'app-notas-list',
  standalone: true,
  imports: [
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatChipsModule,
    LoadingSpinnerComponent,
  ],
  templateUrl: './notas-list.html',
  styleUrl: './notas-list.scss',
})
export class NotasListComponent implements OnInit {
  private readonly notaService = inject(NotaService);
  private readonly router = inject(Router);

  notas: Nota[] = [];
  carregando = false;
  colunas = ['num_seq', 'status', 'acoes'];

  ngOnInit(): void {
    this.carregarNotas();
  }

  carregarNotas(): void {
    this.carregando = true;
    this.notaService.listar().subscribe({
      next: (notas) => {
        this.notas = notas;
        this.carregando = false;
      },
      error: () => {
        this.carregando = false;
      },
    });
  }

  novaNota(): void {
    this.carregando = true;
    this.notaService.criar().subscribe({
      next: (nota) => {
        // Já navega direto para o detalhe da nota recém criada
        this.router.navigate(['/notas', nota.id]);
      },
      error: () => {
        this.carregando = false;
      },
    });
  }

  verDetalhe(id: string): void {
    this.router.navigate(['/notas', id]);
  }
}
