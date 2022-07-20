import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-users-list',
  templateUrl: './users-list.page.html',
  styleUrls: ['./users-list.page.scss'],
})
export class UsersListPage implements OnInit {

  constructor() { }

  ngOnInit() {
  }

  users = [
    {
      name: "Frank tchatseu",
      groupe: "Service client",
      phone: "696812610",
      avatar:"https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"
    },
    {
      name: "Yannick Youmie",
      groupe: "Service client",
      phone: "696812610",
      avatar:"/assets/img/speakers/rabbit.jpg"
    },
    {
      name: "Prestilien Pindoh",
      groupe: "Service client",
      phone: "696812610",
      avatar:"/assets/img/speakers/puppy.jpg"
    },
  ];

  getUsers() {
    
  }

  deleteuser () {
    
  }
}
