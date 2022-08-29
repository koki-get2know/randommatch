import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { FormsModule } from "@angular/forms";

import { IonicModule } from "@ionic/angular";

import { UsersListPageRoutingModule } from "./users-list-routing.module";

import { UsersListPage } from "./users-list.page";
import { SharedModule } from "../../../shared/shared.module";

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    UsersListPageRoutingModule,
    SharedModule,
  ],
  declarations: [UsersListPage],
})
export class UsersListPageModule {}
