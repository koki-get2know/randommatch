import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { IonicModule } from '@ionic/angular';

import { MatchingGroupPageRoutingModule } from './matching-group-routing.module';

import { MatchingGroupPage } from './matching-group.page';
import { IonicSelectableModule } from 'ionic-selectable';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    MatchingGroupPageRoutingModule,
    IonicSelectableModule,
    ReactiveFormsModule
  ],
  declarations: [MatchingGroupPage]
})
export class MatchingGroupPageModule {}
