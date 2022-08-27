import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { FormsModule } from "@angular/forms";

import { IonicModule } from "@ionic/angular";

import { MatchPageRoutingModule } from "./matching-routing.module";

import { MatchingPage } from "./matching.page";
import { SharedModule } from "../../shared/shared.module";
import { MatchingSimplePageModule } from "../matching-simple/matching-simple.module";
import { MatchingGroupPageModule } from "../matching-group/matching-group.module";

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    MatchPageRoutingModule,
    MatchingSimplePageModule,
    MatchingGroupPageModule,
    SharedModule,
  ],
  declarations: [MatchingPage],
})
export class MatchingPageModule {}
