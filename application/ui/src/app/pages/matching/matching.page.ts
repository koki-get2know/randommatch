import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { UsersService } from '../../services/users.service';
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
      nomber:6,
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
      nomber:6,
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
  userSelected = [];

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
      this.userSelected.push( user );
      console.log( this.userSelected );
    }
    else {
      this.onRemoveUserSelected( user.id );
      console.log( this.userSelected );
    }

  }
  // when user is unchecked, it should be remove
  onRemoveUserSelected(id: number) {
    let index = this.userSelected.findIndex(d => d.id === id); //find index in your array
    this.userSelected.splice(index, 1);
    event.stopPropagation();
}
  // select a group of user
  selectGroup(event, group){
  
  this.userSelected = [];
  if ( event.target.checked == false ) {
      this.userSelected = group.users;
      console.log( this.userSelected );
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
    const formData = new FormData();
    const datamatch = {
      matchingsize: this.form.matchingsize.value,
      usersSelected: this.userSelected
    };
    //formData.append( 'matching_size', this.form.matchingsize.value );
    this.matchService.makematch(datamatch)
      .then( resp => {
        console.log(datamatch);
        console.log( "Matching effectué avec success" );
        this.initForm();
        // show the result
        this.matchingresult(datamatch);
      })
      .catch( err => {
        console.log( err );
        
      })
  }

  // matching result
  matchingresult(datamatch) {
    this.router.navigate(['/matching-result',{'data': datamatch} ]);
  }
}
