import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { MatchingGroupPage } from './matching-group.page';

const routes: Routes = [
  {
    path: '',
    component: MatchingGroupPage
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class MatchingGroupPageRoutingModule {}
