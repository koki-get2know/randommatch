import { Component, OnInit } from '@angular/core';
import { UsersService } from '../../../services/users.service';

@Component({
  selector: 'app-users-list',
  templateUrl: './users-list.page.html',
  styleUrls: ['./users-list.page.scss'],
})
export class UsersListPage implements OnInit {

  constructor(private userService:UsersService,) { }

  ngOnInit () {
    console.log( "init" );
    this.getUsers();
  }

  usersgroups = [
    {
      groupe: "Groupe A",
      nomber:6,
      users: [
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
      ]
    },
    {
      groupe: "Groupe B",
      nomber:6,
      users: [
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
      ]
    },
  ];

  uploadCsv ( fileChangeEvent ) {
    
  const photo = fileChangeEvent.target.files[ 0 ];
    
  let formData = new FormData();
  formData.append("photo", photo, photo.name);
  console.log( photo.name );
  this.userService.uploadCsv( formData )
    .then( resp => {
      console.log( resp );
      this.getUsers();
        
      } )
    .catch( err => {
      console.log( err );
    });
    
 }
  getUsers() {
    this.userService.get()
    .then( resp => {
      console.log( resp );
      //this.usersgroups = resp;
      } )
    .catch( err => {
      console.log( err );
    });
  }

}
