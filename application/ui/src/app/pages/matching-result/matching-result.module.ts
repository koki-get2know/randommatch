import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { FormsModule } from "@angular/forms";

import { IonicModule } from "@ionic/angular";

import { MatchingResultPageRoutingModule } from "./matching-result-routing.module";

import { MatchingResultPage } from "./matching-result.page";
import { SharedModule } from "../../shared/shared.module";

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    MatchingResultPageRoutingModule,
    SharedModule,
  ],
  declarations: [MatchingResultPage],
})
export class MatchingResultPageModule {}
