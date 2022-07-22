import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { MatchingResultPage } from './matching-result.page';

const routes: Routes = [
  {
    path: '',
    component: MatchingResultPage
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class MatchingResultPageRoutingModule {}
