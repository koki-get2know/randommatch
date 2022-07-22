import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Matching } from '../../services/users.service';

@Component({
  selector: 'app-matching-result',
  templateUrl: './matching-result.page.html',
  styleUrls: ['./matching-result.page.scss'],
})
export class MatchingResultPage implements OnInit {

  constructor(private route: ActivatedRoute) { }

  ngOnInit () {
    const param: Matching[] = this.route.snapshot.params[ 'data' ];
    console.log( param );
  }

  usersgroups = [
    {
      group: "Match 1",
      number: 6,
      users: [
        {
        name: "Frank tchatseu",
        group: "Service client",
        phone: "696812610",
        avatar:"https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"
      },
      {
        name: "Yannick Youmie",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/rabbit.jpg"
      },
      {
        name: "Prestilien Pindoh",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/puppy.jpg"
      },
      ]
    },
    {
      group: "Match 2",
      number: 6,
      users: [
        {
        name: "Frank tchatseu",
        group: "Service client",
        phone: "696812610",
        avatar:"https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"
      },
      {
        name: "Yannick Youmie",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/rabbit.jpg"
      },
      {
        name: "Prestilien Pindoh",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/puppy.jpg"
      },
      ]
    },
    {
      group: "Match 3",
      number:6,
      users: [
        {
        name: "Frank tchatseu",
        group: "Service client",
        phone: "696812610",
        avatar:"https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"
      },
      {
        name: "Yannick Youmie",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/rabbit.jpg"
      },
      {
        name: "Prestilien Pindoh",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/puppy.jpg"
      },
      ]
    },
    {
      group: "Match 4",
      number:6,
      users: [
        {
        name: "Frank tchatseu",
        group: "Service client",
        phone: "696812610",
        avatar:"https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"
      },
      {
        name: "Yannick Youmie",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/rabbit.jpg"
      },
      {
        name: "Prestilien Pindoh",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/puppy.jpg"
      },
      ]
    },
  ];


}
