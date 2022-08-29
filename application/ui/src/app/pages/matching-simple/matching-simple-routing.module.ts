import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";

import { MatchingSimplePage } from "./matching-simple.page";

const routes: Routes = [
  {
    path: "",
    component: MatchingSimplePage,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class MatchingSimplePageRoutingModule {}
