import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";

import { IonicModule } from "@ionic/angular";

import { MatchingGroupPageRoutingModule } from "./matching-group-routing.module";

import { MatchingGroupPage } from "./matching-group.page";
import { IonicSelectableModule } from "ionic-selectable";
import { SharedModule } from "../../shared/shared.module";
import { SwiperModule } from "swiper/angular";

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    MatchingGroupPageRoutingModule,
    IonicSelectableModule,
    ReactiveFormsModule,
    SwiperModule,
    SharedModule,
  ],
  declarations: [MatchingGroupPage],
})
export class MatchingGroupPageModule {}
