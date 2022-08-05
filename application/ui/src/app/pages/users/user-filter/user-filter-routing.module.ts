import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { UserFilterPage } from './user-filter.page';

const routes: Routes = [
  {
    path: '',
    component: UserFilterPage
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class UserFilterPageRoutingModule {}
