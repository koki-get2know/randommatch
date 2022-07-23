import { Component, OnInit } from '@angular/core';
import { Matching } from '../../services/users.service';

@Component({
  selector: 'app-matching-result',
  templateUrl: './matching-result.page.html',
  styleUrls: ['./matching-result.page.scss'],
})
export class MatchingResultPage implements OnInit {

  matchings: Matching[] = [];

  constructor() { 
  }

  ngOnInit () {
    this.matchings = history.state.matchings;
  }
  




}
