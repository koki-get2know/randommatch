import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { IonicModule } from "@ionic/angular";

import { AccountPage } from "./account";
import { AccountPageRoutingModule } from "./account-routing.module";
import { SharedModule } from "../../shared/shared.module";

@NgModule({
  imports: [CommonModule, IonicModule, AccountPageRoutingModule, SharedModule],
  declarations: [AccountPage],
})
export class AccountModule {}
