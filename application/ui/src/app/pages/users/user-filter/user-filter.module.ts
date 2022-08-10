import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { IonicModule } from '@ionic/angular';

import { UserFilterPageRoutingModule } from './user-filter-routing.module';

import { UserFilterPage } from './user-filter.page';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    UserFilterPageRoutingModule
  ],
  declarations: [UserFilterPage]
})
export class UserFilterPageModule {}
