import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";

import { IonicModule } from "@ionic/angular";

import { MatchingSimplePageRoutingModule } from "./matching-simple-routing.module";

import { MatchingSimplePage } from "./matching-simple.page";
import { IonicSelectableModule } from "ionic-selectable";
import { SwiperModule } from "swiper/angular";
import { SharedModule } from "../../shared/shared.module";

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    ReactiveFormsModule,
    MatchingSimplePageRoutingModule,
    IonicSelectableModule,
    SwiperModule,
    SharedModule,
  ],
  declarations: [MatchingSimplePage],
})
export class MatchingSimplePageModule {}
