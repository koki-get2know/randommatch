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

  selected_forbidden_connexion: [];
  userstoforbidden =[];
  usersconnexionforbidden =[];

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
      if ( this.usersconnexionforbidden.length === 0 ) {
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
        if ( diffUser.length === 0 ) {
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
      const users = [];
      const randomgroup = `${ lorem.generateWords( 2 ) } `;
      for (let i=1; i<usersNumber; i++) {
        const user = {
          id: Math.floor(Math.random() * Date.now()),
          name: lorem.generateWords(2),
          group: randomgroup,
          avatar: this.matchService.generateAvatarSvg()
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
    const index = this.usersSelected.findIndex(d => d.id === id); //find index in your array
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
    

    const users: User[] = [];
    const forbiddenConnections: User[][] = [];
    for (const selected of this.usersSelected)
    {
      users.push({userId: selected.name, avatar: selected.avatar})
    }
    for (const connection of this.usersconnexionforbidden) {
      const newConnection = [];
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
      .subscribe( ( matchings: Matching[] ) => {
        console.log(matchings);
        matchings.forEach(match => match.users.forEach(user => {
          user.avatar = matchingRequest.users.find(usr => usr.userId === user.userId)?.avatar;
        }));
        
        this.matchingresult(matchings);
      })
  }

  // matching result
  matchingresult(matchings: Matching[]) {
    const navigationExtras: NavigationExtras = {
      state: {
        matchings
      }
    };
    this.router.navigate(['/matching-result'],navigationExtras);
  }
  
}
