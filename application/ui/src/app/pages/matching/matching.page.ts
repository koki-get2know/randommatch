import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { UsersService, MatchingReq, User, Matching } from '../../services/users.service';
import { NavController } from '@ionic/angular';
import { NavigationExtras, Router } from '@angular/router';
import { LoremIpsum } from 'lorem-ipsum';
import { IonicSelectableComponent } from 'ionic-selectable';
@Component( {
  selector: 'app-matching',
  templateUrl: './matching.page.html',
  styleUrls: ['./matching.page.scss'],
} )

export class MatchingPage implements OnInit {
  
  matchingForm: FormGroup;
  checked: any;
  usersgroups = [];
   
  usersSelected = [];

  isLoading = false;
  isError = false;
  isSuccess = false;
  isSubmitted = false;
  selected_forbidden_connexion: [];
  userstoforbidden =[];
  usersconnexionforbidden: User[][]=[];

  avatars = ["/assets/img/speakers/bear.jpg", "/assets/img/speakers/cheetah.jpg", "/assets/img/speakers/duck.jpg", 
  "/assets/img/speakers/eagle.jpg", "/assets/img/speakers/elephant.jpg", "/assets/img/speakers/giraffe.jpg", 
  "/assets/img/speakers/iguana.jpg", "/assets/img/speakers/kitten.jpg", "/assets/img/speakers/lion.jpg",
  "/assets/img/speakers/mouse.jpg", "/assets/img/speakers/puppy.jpg", "/assets/img/speakers/rabbit.jpg",
   "/assets/img/speakers/turtle.jpg",
   "https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"];

  @ViewChild('selectComponent') selectComponent:IonicSelectableComponent
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
  clear() {
    this.selectComponent.clear();
    this.selectComponent.close();
    this.selected_forbidden_connexion = [];
    
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
        this.userstoforbidden.push(user);
      }

      let group = {
        group: randomgroup,
        users: users 
      };
      
      usersgroups.push( group );
      
    }
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
