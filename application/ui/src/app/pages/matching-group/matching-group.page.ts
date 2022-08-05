import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { UsersService, MatchingReq, User, Matching } from '../../services/users.service';
import { NavController } from '@ionic/angular';
import { NavigationExtras, Router } from '@angular/router';
import { LoremIpsum } from 'lorem-ipsum';
import { IonicSelectableComponent } from 'ionic-selectable';
@Component({
  selector: 'app-matching-group',
  templateUrl: './matching-group.page.html',
  styleUrls: ['./matching-group.page.scss'],
})
export class MatchingGroupPage implements OnInit {

  matchingForm: FormGroup;
  checked: any;
  usersgroups = [];
   
  usersSelected = [];


  userslist = [
    {
      "UserId": "",
      "Name": "Pins Prestilien",
      "Email": "pinsdev24@gmail.com",
      "Groups": [
        ""
      ],
      "Genre": "Male",
      "Birthday": "10/01",
      "Hobbies": [
        "Data science",
        "Space",
        "Télévison"
      ],
      "MatchPreference": [
        "girls"
      ],
      "MatchPreferenceTime": [
        "14:00PM"
      ],
      "PositionHeld": "CEO",
      "MultiMatch": false,
      "PhoneNumber": "699999999",
      "Departement": "Informatique",
      "Location": "Kao",
      "Seniority": "",
      "Role": "super-user",
      "NumberOfMatching": 0,
      "NumberMatchingAccepted": 0,
      "NumberMatchingDeclined": 0,
      "AverageMatchingRate": 0
    },
    {
      "UserId": "",
      "Name": "Frank Tchatseu",
      "Email": "pinsdev24@gmail.com",
      "Groups": [
        "DS",
        "IA",
        "SPACE"
      ],
      "Genre": "Male",
      "Birthday": "10/01",
      "Hobbies": [
        "Jeux vidéos",
        "Musique"
      ],
      "MatchPreference": [
        "same groups"
      ],
      "MatchPreferenceTime": [
        "14:00PM"
      ],
      "PositionHeld": "Admin",
      "MultiMatch": false,
      "PhoneNumber": "699999999",
      "Departement": "Math",
      "Location": "Fpol",
      "Seniority": "",
      "Role": "user",
      "NumberOfMatching": 0,
      "NumberMatchingAccepted": 0,
      "NumberMatchingDeclined": 0,
      "AverageMatchingRate": 0
    },
    {
      "UserId": "",
      "Name": "Pins Prestilien",
      "Email": "pinsdev24@gmail.com",
      "Groups": [
        "Foot"
      ],
      "Genre": "Male",
      "Birthday": "10/01",
      "Hobbies": [
        "Dance"
      ],
      "MatchPreference": [
        "pretty girl"
      ],
      "MatchPreferenceTime": [
        "14:00PM"
      ],
      "PositionHeld": "Chef",
      "MultiMatch": false,
      "PhoneNumber": "699999999",
      "Departement": "Biology",
      "Location": "Soa",
      "Seniority": "",
      "Role": "user",
      "NumberOfMatching": 0,
      "NumberMatchingAccepted": 0,
      "NumberMatchingDeclined": 0,
      "AverageMatchingRate": 0
    }
  ];

  ColorsTags = [
    "twitter",
    "instagram",
    "dark"
  ]
  
  getRandomColor () {
    const min = 0;
    const max = 2;
    const index = Math.floor( Math.random() * ( max - min + 1 ) ) + min;
    return this.ColorsTags[3%(index+1)];
  }


  isLoading = false;
  isError = false;
  isSuccess = false;
  isSubmitted = false;
  selected_forbidden_connexion: [];
  userstoforbidden =[];
  usersconnexionforbidden: User[][] = [];

  users_toselect_group1 = [];
  users_toselect_group2 = [];
  users_selected_group1 = [];
  users_selected_group2 = [];
  result_selected_group1 = [];
  result_selected_group2= [];

  avatars = ["/assets/img/speakers/bear.jpg", "/assets/img/speakers/cheetah.jpg", "/assets/img/speakers/duck.jpg", 
  "/assets/img/speakers/eagle.jpg", "/assets/img/speakers/elephant.jpg", "/assets/img/speakers/giraffe.jpg", 
  "/assets/img/speakers/iguana.jpg", "/assets/img/speakers/kitten.jpg", "/assets/img/speakers/lion.jpg",
  "/assets/img/speakers/mouse.jpg", "/assets/img/speakers/puppy.jpg", "/assets/img/speakers/rabbit.jpg",
   "/assets/img/speakers/turtle.jpg",
   "https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"];

  @ViewChild( 'selectComponent' ) selectComponent: IonicSelectableComponent
  @ViewChild( 'selectComponentGroup1' ) selectComponentGroup1: IonicSelectableComponent
  @ViewChild( 'selectComponentGroup2') selectComponentGroup2:IonicSelectableComponent
  constructor(private formBuilder: FormBuilder,private matchService:UsersService,
    public navCtrl: NavController, private router: Router ) { 
    
  }

  ngOnInit () {
    this.usersgroups = this.generateUsers();
    this.initForm();
  }

  initForm() {
    this.matchingForm = this.formBuilder.group({
      matchingsize: ['', Validators.required],
    });
  }

  portChange(event: {
    component: IonicSelectableComponent,
    value: any
  } ) {
    console.log( "Selec" );
    if ( this.selected_forbidden_connexion.length > 0 ) {
      this.usersconnexionforbidden.push( this.selected_forbidden_connexion );
    }
    console.log( this.usersconnexionforbidden );
    console.log( 'port:', event.value );

    this.clear();
  }
  portChangeGroup1(event: {
    component: IonicSelectableComponent,
    value: any
  } ) {
    console.log( "Selec group 1" );
   
    console.log( this.users_selected_group1 );
    this.result_selected_group1 = this.users_selected_group1;
    // users to forbidden must be selected among the group 1 and 2
    this.userstoforbidden = [];
    this.userstoforbidden = this.userstoforbidden.concat( this.result_selected_group1 );
    this.userstoforbidden = this.userstoforbidden.concat(this.result_selected_group2);

    const forbid = 

    console.log( "LA CONCATs" );
    console.log( this.userstoforbidden );
    this.clearGroup1();
  }
  portChangeGroup2(event: {
    component: IonicSelectableComponent,
    value: any
  } ) {
    console.log( "Select group 2" );
    this.result_selected_group2 = this.users_selected_group2;

    this.userstoforbidden = this.userstoforbidden.concat(this.result_selected_group2);
    
    this.clearGroup2();
  }
  clear() {
    this.selectComponent.clear();
    this.selectComponent.close();
    this.selected_forbidden_connexion = [];
    
  }
  clearGroup1() {
    this.selectComponentGroup1.clear();
    this.selectComponentGroup1.close();
  }
  clearGroup2() {
    this.selectComponentGroup2.clear();
    this.selectComponentGroup2.close();
    
  }
  addforbiddenUsersItem() {
    this.selectComponent.confirm ();
    this.selectComponent.close(); 
    
  }
  generateUsers() {
    const lorem = new LoremIpsum({
      sentencesPerParagraph: {
        max: 8,
        min: 4
      },
      wordsPerSentence: {
        max: 16,
        min: 4
      }
    });
    const min = 6;
    const max = 30;
    const usersNumber = Math.floor(Math.random() * (max - min + 1)) + min;

    let usersgroups = [];
    for (let g=1; g < 3; g++) {
      let users = [];
      let randomgroup = `Group ${ lorem.generateWords( 2 ) }`;
      for (let i=1; i<usersNumber; i++) {
        let avatarId = Math.floor(Math.random() * (this.avatars.length));
        let user = {
          id: i,
          name: lorem.generateWords(2),
          group: randomgroup,
          avatar: this.avatars[avatarId]
        };
        users.push( user );
        //this.userstoforbidden.push( user );
        this.users_toselect_group1.push( user );
        this.users_toselect_group2.push( user );

      }

      let group = {
        group: randomgroup,
        users: users 
      };
      
      usersgroups.push( group );
      
    }
    console.log( this.userstoforbidden );
    return usersgroups;
}

  get form() {
    return this.matchingForm.controls;
  }

  selectUsers(event,user) {
  
    if (!!event.target.checked === false ) {
      this.usersSelected.push( user );
    }
    else {
      this.onRemoveusersSelected( user.id );
    }

  }
  // select users per groups

  selectUsersGroup1(event,user) {
  
    if ( !!event.target.checked === false ) {
      this.users_selected_group1.push( user );
    }
    else {
      this.onRemoveusersSelectedGroup1( user.id );
    }
  }

  // select users per groups 2

  selectUsersGroup2(event,user) {
  
    if ( !!event.target.checked === false ) {
      this.users_selected_group2.push( user );
    }
    else {
      this.onRemoveusersSelectedGroup2( user.id );
    }
  }


  // when user is unchecked, it should be remove
  onRemoveusersSelected(id: number) {
    let index = this.usersSelected.findIndex(d => d.id === id); //find index in your array
    this.usersSelected.splice(index, 1);
    event.stopPropagation();
  }
  // when user is unchecked, it should be remove
  onRemoveusersSelectedGroup1(id: number) {
    let index = this.users_selected_group1.findIndex(d => d.id === id); //find index in your array
    this.users_selected_group1.splice(index, 1);
    event.stopPropagation();
  }
  // when user is unchecked, it should be remove
  onRemoveusersSelectedGroup2(id: number) {
    let index = this.users_selected_group2.findIndex(d => d.id === id); //find index in your array
    this.users_selected_group2.splice(index, 1);
    event.stopPropagation();
 }
  // select a group of user
  selectGroup(event, group){
  
  this.usersSelected = [];
  if ( event.target.checked == false ) {
      this.usersSelected = group.users;
    }
  }
  
  onSubmit() {
    this.ramdommatch();
  }

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
      users,
      forbiddenConnections: this.usersconnexionforbidden

    };

    this.matchService.makematch(matchingRequest)
      .subscribe( ( res: Matching[] ) => {
        console.log( matchingRequest );
        this.matchingresult(res);
      })
  }

  // matching result
  matchingresult(matchings: Matching[]) {
    for(const matching of matchings) {
      for (const user of matching.users) {
        let avatarId = Math.floor(Math.random() * (this.avatars.length));
        user.avatar = this.avatars[avatarId];
      }
    }
    const navigationExtras: NavigationExtras = {
      state: {
        matchings
      }
    };
    this.router.navigate(['/matching-result'],navigationExtras);
  }

}
