import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { UsersService, MatchingReq, User, Matching, MatchingGroupReq } from '../../services/users.service';
import { NavController, ToastController } from '@ionic/angular';
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


  userslist = [];

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
  groups: User[][] = [];

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
    public navCtrl: NavController, private router: Router,public toastController: ToastController ) { 
    
  }

  ngOnInit () {
    //this.usersgroups = this.generateUsers();
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

 
    portChange(event: {
    component: IonicSelectableComponent,
    value: any} ) {
    console.log( "Selec" );
    // just add if the list in not empty
    if ( this.selected_forbidden_connexion.length > 1 ) {
      if ( this.usersconnexionforbidden.length == 0 ) {
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
    this.userstoforbidden = this.userstoforbidden.concat( this.result_selected_group2 );
    // in the second group, just keep all the users who have not been selected in the group 1
    this.users_toselect_group2 = this.compareconnection( this.users_toselect_group1, this.result_selected_group1 );

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
    
    // in the second group, just keep all the users who have not been selected in the group 1
    this.users_toselect_group1 = this.compareconnection( this.users_toselect_group2, this.result_selected_group2 );

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
  compareconnection<User>(forbconnec1:any,forbconnec2:any) {
    return forbconnec1.filter((element) => {
        return !forbconnec2.some(elt2 => element.id === elt2.id);
      });
    
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
      users.push({id: selected.name})
    }
    this.groups.push( this.result_selected_group1 );
    this.groups.push( this.result_selected_group2 );
    const matchingRequest: MatchingGroupReq = {
      size: Number(this.form.matchingsize.value),
      groups:this.groups,
      forbiddenConnections: this.usersconnexionforbidden

    };

    console.log( "Data Match Group" );
    console.log( matchingRequest );

    this.matchService.makematchgroup(matchingRequest)
      .subscribe( ( res: Matching[] ) => {
        if ( res!==null ) {
          console.log( matchingRequest );
          this.matchingresult(res);
        }
        else {
          this.presentToast("No match generated!");
        }
       
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

   async presentToast(message) {
    const toast = await this.toastController.create({
      message: message,
      duration: 2000
    });
    toast.present();
  }
}
