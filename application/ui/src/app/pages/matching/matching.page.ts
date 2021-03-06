import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { UsersService, MatchingReq, User, Matching } from '../../services/users.service';
import { NavController, ToastController } from '@ionic/angular';
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
  usersconnexionforbidden =[];

  avatars = ["/assets/img/speakers/bear.jpg", "/assets/img/speakers/cheetah.jpg", "/assets/img/speakers/duck.jpg", 
  "/assets/img/speakers/eagle.jpg", "/assets/img/speakers/elephant.jpg", "/assets/img/speakers/giraffe.jpg", 
  "/assets/img/speakers/iguana.jpg", "/assets/img/speakers/kitten.jpg", "/assets/img/speakers/lion.jpg",
  "/assets/img/speakers/mouse.jpg", "/assets/img/speakers/puppy.jpg", "/assets/img/speakers/rabbit.jpg",
   "/assets/img/speakers/turtle.jpg",
   "https://avatars.githubusercontent.com/u/50463560?s=400&u=d082fa7694a0d14dc2e464adc8e6e7ef4ce49aaa&v=4"];

  @ViewChild('selectComponent') selectComponent:IonicSelectableComponent
  constructor(private formBuilder: FormBuilder,private matchService:UsersService,
    public navCtrl: NavController, private router: Router,
    public toastController: ToastController) { 
    
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
  async presentToast(message) {
    const toast = await this.toastController.create({
      message: message,
      duration: 2000
    });
    toast.present();
  }

  userChange(event: {
    component: IonicSelectableComponent,
    value: any} ) {
    // just add if the list in not empty
    if ( this.selected_forbidden_connexion.length > 1 ) {
      if ( this.usersconnexionforbidden.length == 0 ) {
        this.usersconnexionforbidden.push( this.selected_forbidden_connexion );
      }
      else {
        if ( !this.forbiddenConnectionAlreadyExist( this.selected_forbidden_connexion ) ) {
          console.log( "Unexisting link" );
          this.usersconnexionforbidden.push( this.selected_forbidden_connexion );
        }
        else {
          this.presentToast("this connection already exist!");
        }
      }
    } else {
      this.presentToast("Please select more than one user!");
    }
    this.clear();
  }
  forbiddenConnectionAlreadyExist ( newconnection: User[] ): boolean {
    let i = 0;
    while ( i < this.usersconnexionforbidden.length ) {
      let element = this.usersconnexionforbidden[i];
      if ( element.length == newconnection.length ) {
        const diffUser = this.compareconnection( element, newconnection );
        if ( diffUser.length == 0 ) {
          console.log( diffUser.length);
          return true;
        }
      }
      i++;
    }
    return false;
  }

  compareconnection(forbconnec1:any,forbconnec2:any) {
    return forbconnec1.filter((element) => {
        return !forbconnec2.some(elt2 => element.id === elt2.id);
      });
    
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

  removeConnection (index) {
    this.usersconnexionforbidden.splice(index, 1);
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
      let randomgroup = `${ lorem.generateWords( 2 ) } `;
      for (let i=1; i<usersNumber; i++) {
        let avatarId = Math.floor( Math.random() * ( this.avatars.length ) );
        // generate unique id
        let user = {
          id: Math.floor(Math.random() * Date.now()),
          name: lorem.generateWords(2),
          group: randomgroup,
          avatar: this.avatars[avatarId]
        };
        users.push( user );
        this.userstoforbidden.push(user);
      }

      let group = {
        group: `Group ${g} ${randomgroup}`,
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
    let forbiddenConnections: User[][] = [];
    for (let selected of this.usersSelected)
    {
      users.push({userId: selected.name})
    }
    for (let connection of this.usersconnexionforbidden) {
      let newConnection = [];
      for (let item of connection) {
        newConnection.push({userId: item.name});
      }
      forbiddenConnections.push(newConnection);
    }
    const matchingRequest: MatchingReq = {
      size: Number(this.form.matchingsize.value),
      users,
      forbiddenConnections
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
