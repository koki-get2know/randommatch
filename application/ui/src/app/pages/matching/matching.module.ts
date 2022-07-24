import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule,ReactiveFormsModule } from '@angular/forms';

import { IonicModule } from '@ionic/angular';

import { MatchingPageRoutingModule } from './matching-routing.module';

import { MatchingPage } from './matching.page';
import { IonicSelectableModule } from 'ionic-selectable';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    ReactiveFormsModule,
    MatchingPageRoutingModule,
    IonicSelectableModule
    
  ],
  declarations: [MatchingPage]
})
export class MatchingPageModule {}
