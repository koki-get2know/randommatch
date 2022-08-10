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
  usersconnexionforbidden: User[][] = [];
  
  ColorsTags = [
    "twitter",
    "instagram",
    "dark"
  ]

  @ViewChild('selectComponent') selectComponent:IonicSelectableComponent
  constructor(private formBuilder: FormBuilder,private matchService:UsersService,
    public navCtrl: NavController, private router: Router,
    public toastController: ToastController) { 
    
  }

  ngOnInit () {
    const storagevalue= localStorage.getItem( "userlist" );
    this.usersgroups = storagevalue ? JSON.parse( storagevalue ) : [];
    this.initusermatch();
    this.initForm();
  }

  isIndeterminate:boolean;
  masterCheck:boolean;
  checkBoxList: any;
  
  checkMaster ( event ) {
    this.usersSelected = [];
    setTimeout( () => {
      if ( this.masterCheck ) {
        this.usersgroups.forEach(user => {
        user.isChecked = this.masterCheck;
        this.usersSelected.push( user );
        });
      }
      else {
        this.usersgroups.forEach(user => {
        user.isChecked = this.masterCheck;
        this.onRemoveusersSelected( user.id );
        });
      }
      
    });
  }

  initusermatch () {
    this.usersgroups.forEach(obj => {
        obj.isChecked = false;
       // obj = { ...obj, isChecked: false};
      });
  }
 selectUsers(event: PointerEvent,user) {
    if ((event.target as HTMLInputElement).checked === false ) {
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
  }
    checkEvent(event: PointerEvent,user) {
    const totalItems = this.usersgroups.length;
    let checked = 0;
    
    if ((event.target as HTMLInputElement).checked === false ) {
      this.usersSelected.push( user );
      checked++;
      
    }
    else {
      this.onRemoveusersSelected( user.id );
      checked--;
    }
    if (checked > 0 && checked < totalItems) {
      //If even one item is checked but not all
      this.isIndeterminate = true;
      this.masterCheck = false;
    } else if (checked === totalItems) {
      //If all are checked
      this.masterCheck = true;
      this.isIndeterminate = false;
    } else {
      //If none is checked
      this.isIndeterminate = false;
      this.masterCheck = false;
    }
  }

  getRandomColor () {
    const min = 0;
    const max = 2;
    const index = Math.floor( Math.random() * ( max - min + 1 ) ) + min;
    return this.ColorsTags[3%(index+1)];
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
          console.log( "Lien inexistant" );
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
      if ( element.length === newconnection.length ) {
        const diffUser = this.compareconnection( element, newconnection );
        if ( diffUser.length === 0 ) {
          console.log( diffUser.length);
          return true;
        }
      }
      i++;
    }
    return false;
  }

  compareconnection<User>(forbconnec1:any,forbconnec2:any) {
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

  get form() {
    return this.matchingForm.controls;
  }

 
  // select a group of user
  selectGroup(event, group){
  
  this.usersSelected = [];
  if ( event.target.checked === false ) {
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
      users.push({id: selected.id, name: selected.name, avatar: selected.avatar})
    }
    for (const connection of this.usersconnexionforbidden) {
      const newConnection = [];
      for (let item of connection) {
        newConnection.push({id: item.id, name: item.name});
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
        if ( matchings!==null ) {
            console.log(matchings);
            matchings.forEach(match => match.users.forEach(user => {
              user.avatar = matchingRequest.users.find(usr => usr.id === user.id)?.avatar;
            }));
            
            this.matchingresult(matchings);
        }
        else {
          this.presentToast("No match generated!");
        }
        
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
