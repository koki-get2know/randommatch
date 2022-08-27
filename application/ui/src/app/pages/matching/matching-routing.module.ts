import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";

import { MatchingPage } from "./matching.page";

const routes: Routes = [
  {
    path: "",
    children: [
      {
        path: "simple",
        children: [
          {
            path: "",
            loadChildren: () =>
              import("../matching-simple/matching-simple.module").then(
                (m) => m.MatchingSimplePageModule
              ),
          },
        ],
      },
      {
        path: "group",
        children: [
          {
            path: "",
            loadChildren: () =>
              import("../matching-group/matching-group.module").then(
                (m) => m.MatchingGroupPageModule
              ),
          },
        ],
      },
      {
        path: "",
        pathMatch: "full",
        component: MatchingPage,
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class MatchPageRoutingModule {}
