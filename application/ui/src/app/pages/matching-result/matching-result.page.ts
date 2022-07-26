import { Component, OnInit } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { Matching } from '../../services/users.service';

@Component({
  selector: 'app-matching-result',
  templateUrl: './matching-result.page.html',
  styleUrls: ['./matching-result.page.scss'],
})
export class MatchingResultPage implements OnInit {

  matchings: Matching[] = [];

  constructor(private sanitizer: DomSanitizer) { 
  }

  ngOnInit () {
    this.matchings = history.state.matchings;
    this.matchings?.forEach(match => match.users.forEach(user => {
      if (user.avatar) {
        user.avatar = this.sanitizer.bypassSecurityTrustHtml(user.avatar['changingThisBreaksApplicationSecurity']);
      }
    }));
  }
  




}


