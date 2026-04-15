import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';

import { ProdutoService } from '../../../core/services/produto.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner';

@Component({
  selector: 'app-produto-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    LoadingSpinnerComponent,
  ],
  templateUrl: './produto-form.html',
  styleUrl: './produto-form.scss',
})
export class ProdutoFormComponent {
  private readonly produtoService = inject(ProdutoService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly cdr = inject(ChangeDetectorRef);

  salvando = false;

  form = this.fb.group({
    codigo: ['', [Validators.required, Validators.minLength(2)]],
    nome: ['', [Validators.required, Validators.minLength(3)]],
    saldo: [0, [Validators.required, Validators.min(0)]],
  });

  salvar(): void {
    if (this.form.invalid) return;

    this.salvando = true;
    this.cdr.detectChanges();
    this.produtoService
      .criar({
        codigo: this.form.value.codigo!,
        nome: this.form.value.nome!,
        saldo: this.form.value.saldo!,
      })
      .subscribe({
        next: () => {
          this.router.navigate(['/produtos']);
        },
        error: () => {
          // Erro já tratado pelo interceptor global
          this.salvando = false;
          this.cdr.detectChanges();
        },
      });
  }

  cancelar(): void {
    this.router.navigate(['/produtos']);
  }

  // Helpers para exibir erros no template
  erroNome(): string {
    const ctrl = this.form.get('nome');
    if (ctrl?.hasError('required')) return 'Nome é obrigatório';
    if (ctrl?.hasError('minlength')) return 'Nome deve ter ao menos 3 caracteres';
    return '';
  }

  erroCodigo(): string {
    const ctrl = this.form.get('codigo');
    if (ctrl?.hasError('required')) return 'Codigo e obrigatorio';
    if (ctrl?.hasError('minlength')) return 'Codigo deve ter ao menos 2 caracteres';
    return '';
  }

  erroSaldo(): string {
    const ctrl = this.form.get('saldo');
    if (ctrl?.hasError('required')) return 'Saldo é obrigatório';
    if (ctrl?.hasError('min')) return 'Saldo não pode ser negativo';
    return '';
  }
}
