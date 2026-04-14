import { Routes } from '@angular/router';

export const NOTAS_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () => import('./notas-list/notas-list').then((m) => m.NotasListComponent),
  },
  {
    path: ':id/itens',
    loadComponent: () => import('./nota-form/nota-form').then((m) => m.NotaFormComponent),
  },
  {
    path: ':id',
    loadComponent: () => import('./nota-detalhe/nota-detalhe').then((m) => m.NotaDetalheComponent),
  },
];
