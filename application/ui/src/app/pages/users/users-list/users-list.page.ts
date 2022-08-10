import { Component, ViewChild, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AlertController, IonList, IonRouterOutlet, LoadingController, ModalController, ToastController, Config } from '@ionic/angular';

import { ScheduleFilterPage } from '../../schedule-filter/schedule-filter';
import { ConferenceData } from '../../../providers/conference-data';
import { UserData } from '../../../providers/user-data';
import { UserFilterPage } from '../user-filter/user-filter.page';
import { UsersService } from '../../../services/users.service';
import { Observable } from 'rxjs';
import { shareReplay } from 'rxjs/operators';

@Component({
  selector: 'app-users-list',
  templateUrl: './users-list.page.html',
  styleUrls: ['./users-list.page.scss'],
})
export class UsersListPage implements OnInit {

  groups: any = [];

  constructor(
    public router: Router,
    public toastCtrl: ToastController,
    private userService: UsersService
  ) { 
    
  }

  userslist = [];
  key: string = 'userlist';
  isloading: boolean = false;
  
  storeUsersList() {
    localStorage.setItem(this.key, JSON.stringify(this.userslist));
    const storagevalue= localStorage.getItem( this.key );
    const val = storagevalue ? JSON.parse(storagevalue) : []
  }
  
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


  ngOnInit() {
    const storagevalue= localStorage.getItem( "userlist" );
    this.userslist = storagevalue ? JSON.parse( storagevalue ) : [];
  }

  tagclick () {
    
  }

  checkResponseUrl = "";
  
  uploadCsv ( event ) {
    this.isloading = true;
    for (const file of event.target.files) {
      const fileData = new FormData();
      fileData.append("file", file);
      this.userService.uploadUsersFile( fileData )
        .subscribe( resp => {
          if ( resp.status === 202 ) {
            this.checkResponseUrl = resp.headers.get( "Location" );
            this.checkserverresponse();
          }
          
        } );
    }
  }
  
  checkserverresponse () {
    let responsestatus = "";
    const limitedInterval = setInterval(() => {
      this.userService.availabilyofusers( this.checkResponseUrl )
        .subscribe( resp => {
          responsestatus = resp[ "status" ];
          if ( responsestatus === 'Done'  ) {
            this.getuserList();
            this.isloading = false;
          } else if (responsestatus !== 'Running') {
            this.isloading = false;
          }
        } );
      
        if ((responsestatus !== 'Running' && responsestatus !== 'Pending')) {
          clearInterval(limitedInterval);
          console.log( 'interval cleared!' );

        }
    }, 300);
  }

  getuserList () {
    this.userService.getUsersdata()
        .subscribe( resp => {
          this.userslist = resp.data;
          this.storeUsersList();
          
        } );
  }
}
