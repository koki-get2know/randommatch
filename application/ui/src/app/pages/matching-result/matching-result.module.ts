import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { IonicModule } from '@ionic/angular';

import { MatchingResultPageRoutingModule } from './matching-result-routing.module';

import { MatchingResultPage } from './matching-result.page';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    MatchingResultPageRoutingModule
  ],
  declarations: [MatchingResultPage]
})
export class MatchingResultPageModule {}
