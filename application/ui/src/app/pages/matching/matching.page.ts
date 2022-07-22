import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { UsersService, MatchingReq, User, Matching } from '../../services/users.service';
import { NavController } from '@ionic/angular';
import { Router } from '@angular/router';
@Component( {
  selector: 'app-matching',
  templateUrl: './matching.page.html',
  styleUrls: ['./matching.page.scss'],
})
export class MatchingPage implements OnInit {

  constructor(private formBuilder: FormBuilder,private matchService:UsersService,public navCtrl: NavController,private router: Router) { }

  ngOnInit () {
    this.initForm();
  }

  createForm: FormGroup;
  checked: any;
  usersgroups = [
    {
      group: "group A",
      number:6,
      users: [
        {
        id:1,
        name: "Frank tchatseu",
        group: "Service client",
        phone: "696812610",
        avatar:"https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"
      },
        {
        id:2,
        name: "Yannick Youmie",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/rabbit.jpg"
      },
        {
        id:3,
        name: "Prestilien Pindoh",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/puppy.jpg"
      },
      ]
    },
    {
      group: "group B",
      number:6,
      users: [
        {
        id:4,
        name: "Frank tchatseu",
        group: "Service client",
        phone: "696812610",
        avatar:"https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"
      },
        {
        id:5,
        name: "Yannick Youmie",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/rabbit.jpg"
      },
        { 
        id:6,
        name: "Prestilien Pindoh",
        group: "Service client",
        phone: "696812610",
        avatar:"/assets/img/speakers/puppy.jpg"
      },
      ]
    },
  ];
  usersSelected = [];

  isLoading = false;
  isError = false;
  isSuccess = false;
  isSubmitted = false;

  initForm() {
    this.createForm = this.formBuilder.group({
      matchingsize: ['', Validators.required],
    });
  }

  get form() {
    return this.createForm.controls;
  }

  selectUsers(event,user) {
  
    if (event.target.checked == false ) {
      this.usersSelected.push( user );
      console.log( this.usersSelected );
    }
    else {
      this.onRemoveusersSelected( user.id );
      console.log( this.usersSelected );
    }

  }
  // when user is unchecked, it should be remove
  onRemoveusersSelected(id: number) {
    let index = this.usersSelected.findIndex(d => d.id === id); //find index in your array
    this.usersSelected.splice(index, 1);
    event.stopPropagation();
}
  // select a group of user
  selectGroup(event, group){
  
  this.usersSelected = [];
  if ( event.target.checked == false ) {
      this.usersSelected = group.users;
    }
  }
  // 
  ramdommatch () {
    
    this.isSubmitted = true;
    this.isError = false;
    this.isSuccess = false;
    this.isLoading = false
    if ( this.form.invalid ) {
      alert( "remplir tous les champs" );
    }
    this.isLoading = true;


    let users: User[] = [];
    for (let selected of this.usersSelected)
    {
      users.push({userId: selected.name})
    }
    const matchingRequest: MatchingReq = {
      size: Number(this.form.matchingsize.value),
      users
    };

    this.matchService.makematch(matchingRequest)
      .subscribe( (res: Matching[]) => {
        console.log(res);
        console.log( "Matching effectué avec success" );
        this.initForm();
        // show the result
        this.matchingresult(res);
      })
  }

  // matching result
  matchingresult(matchings: Matching[]) {
    this.router.navigate(['/matching-result',{'data': matchings} ]);
  }
}
