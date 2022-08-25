import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { IonicModule } from "@ionic/angular";
import { TranslateModule } from "@ngx-translate/core";

@NgModule({
  declarations: [],
  imports: [CommonModule, FormsModule, IonicModule, ReactiveFormsModule],
  exports: [FormsModule, TranslateModule],
})
export class SharedModule {}
