import { Routes } from '@angular/router';

export const PRODUTOS_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./produtos-list/produtos-list').then((m) => m.ProdutosListComponent),
  },
  {
    path: 'novo',
    loadComponent: () => import('./produto-form/produto-form').then((m) => m.ProdutoFormComponent),
  },
];
