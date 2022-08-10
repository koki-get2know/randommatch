import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { UsersService, User, Matching, MatchingGroupReq } from '../../services/users.service';
import { NavController, ToastController } from '@ionic/angular';
import { NavigationExtras, Router } from '@angular/router';
import { IonicSelectableComponent } from 'ionic-selectable';
import { ColorsTags } from '../../constants';
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
  userslist = [];
  ColorsTags = ColorsTags;
  
  isLoading = false;
  isError = false;
  isSuccess = false;
  isSubmitted = false;
  selected_forbidden_connexion: User[];
  userstoforbidden: User[] =[];
  usersconnexionforbidden: User[][] = [];
  groups: User[][] = [];

  users_toselect_group1: User[] = [];
  users_toselect_group2: User[] = [];
  users_selected_group1: User[] = [];
  users_selected_group2: User[] = [];
  result_selected_group1: User[] = [];
  result_selected_group2: User[]= [];


  @ViewChild( 'selectComponent' ) selectComponent: IonicSelectableComponent
  @ViewChild( 'selectComponentGroup1' ) selectComponentGroup1: IonicSelectableComponent
  @ViewChild( 'selectComponentGroup2') selectComponentGroup2:IonicSelectableComponent
  constructor(private formBuilder: FormBuilder,private matchService:UsersService,
    public navCtrl: NavController, private router: Router,public toastController: ToastController ) { 
    
  }

  ngOnInit () {
    const storagevalue= localStorage.getItem( "userlist" );
    this.users_toselect_group1 = storagevalue ? JSON.parse( storagevalue ) : [];
    this.users_toselect_group2 = storagevalue ? JSON.parse( storagevalue ) : [];

    this.initForm();
  }

  initForm() {
    this.matchingForm = this.formBuilder.group({
      matchingsize: ['', Validators.required],
    });
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
          this.presentToast("this connection already exists!");
        }
      }
    } else {
      this.presentToast("Please select more than one user!");
    }
    this.clear();
  }
  
  userChangeGroup1(event: {
    component: IonicSelectableComponent,
    value: any
  } ) {
   
    console.log( this.users_selected_group1 );
    this.result_selected_group1 = this.users_selected_group1;
    // users to forbid must be selected among the group 1 and 2
    this.userstoforbidden = [];
    this.userstoforbidden = this.userstoforbidden.concat( this.result_selected_group1 );
    this.userstoforbidden = this.userstoforbidden.concat( this.result_selected_group2 );
    // in the second group, just keep all the users who have not been selected in the group 1
    this.users_toselect_group2 = this.matchService.compareconnection( this.users_toselect_group1, this.result_selected_group1 );
    this.clearGroup1();
  }
  userChangeGroup2(event: {
    component: IonicSelectableComponent,
    value: any
  } ) {
    this.result_selected_group2 = this.users_selected_group2;

    this.userstoforbidden = this.userstoforbidden.concat(this.result_selected_group2);
    
    // in the second group, just keep all the users who have not been selected in the group 1
    this.users_toselect_group1 = this.matchService.compareconnection( this.users_toselect_group2, this.result_selected_group2 );

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
  forbiddenConnectionAlreadyExist ( newconnection: User[] ): boolean {
    let i = 0;
    while ( i < this.usersconnexionforbidden.length ) {
      let element = this.usersconnexionforbidden[i];
      if ( element.length === newconnection.length ) {
        const diffUser = this.matchService.compareconnection( element, newconnection );
        if ( diffUser.length === 0 ) {
          console.log( diffUser.length);
          return true;
        }
      }
      i++;
    }
    return false;
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

  selectUsersGroup2(event,user:User) {
  
    if ( !!event.target.checked === false ) {
      this.users_selected_group2.push( user );
    }
    else {
      this.onRemoveusersSelectedGroup2( user.id );
    }
  }


  // when user is unchecked, it should be remove
  onRemoveusersSelected(id: string) {
    let index = this.usersSelected.findIndex(d => d.id === id); //find index in your array
    this.usersSelected.splice(index, 1);
  }
  // when user is unchecked, it should be remove
  onRemoveusersSelectedGroup1(id: string) {
    let index = this.users_selected_group1.findIndex(d => d.id === id); //find index in your array
    this.users_selected_group1.splice(index, 1);
  }
  // when user is unchecked, it should be remove
  onRemoveusersSelectedGroup2(id: string) {
    let index = this.users_selected_group2.findIndex(d => d.id === id); //find index in your array
    this.users_selected_group2.splice(index, 1);
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
    if (this.result_selected_group1.length == 0 || this.result_selected_group2.length ==0) {
      return;
    }
    this.isSubmitted = true;
    this.isError = false;
    this.isSuccess = false;
    this.isLoading = false;
    if ( this.form.invalid ) {
      alert( "Fill all the fields" );
    }
    this.isLoading = true;

    this.groups.push( this.result_selected_group1 );
    this.groups.push( this.result_selected_group2 );
    const matchingRequest: MatchingGroupReq = {
      size: Number(this.form.matchingsize.value),
      groups:this.groups,
      forbiddenConnections: this.usersconnexionforbidden

    };

    this.matchService.makematchgroup(matchingRequest)
      .subscribe( ( matchings: Matching[] ) => {
        if ( matchings !== null ) {
          matchings.forEach(match => match.users.forEach(user => {
            for (const users of matchingRequest.groups) {
              for (const usr of users) {
                if (user.id === usr.id ) {
                  user.avatar = usr.avatar;
                  break;
                }
              }
            }
          }));
          this.matchingresult(matchings);
        }
        else {
          this.presentToast("No match generated!");
        }
      });
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

   async presentToast(message) {
    const toast = await this.toastController.create({
      message: message,
      duration: 2000
    });
    toast.present();
  }
}
