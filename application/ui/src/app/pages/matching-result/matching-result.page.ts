import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Matching } from '../../services/users.service';

@Component({
  selector: 'app-matching-result',
  templateUrl: './matching-result.page.html',
  styleUrls: ['./matching-result.page.scss'],
})
export class MatchingResultPage implements OnInit {

  matchings: Matching[] = [];
  avatars = ["/assets/img/speakers/bear.jpg", "/assets/img/speakers/cheetah.jpg", "/assets/img/speakers/duck.jpg", 
  "/assets/img/speakers/eagle.jpg", "/assets/img/speakers/elephant.jpg", "/assets/img/speakers/giraffe.jpg", 
  "/assets/img/speakers/iguana.jpg", "/assets/img/speakers/kitten.jpg", "/assets/img/speakers/lion.jpg",
  "/assets/img/speakers/mouse.jpg", "/assets/img/speakers/puppy.jpg", "/assets/img/speakers/rabbit.jpg",
   "/assets/img/speakers/turtle.jpg"];

  constructor() { 
  }

  ngOnInit () {
    this.matchings = history.state.matchings;
  }



}
