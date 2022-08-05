import { Component, ViewChild, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AlertController, IonList, IonRouterOutlet, LoadingController, ModalController, ToastController, Config } from '@ionic/angular';

import { ScheduleFilterPage } from '../../schedule-filter/schedule-filter';
import { ConferenceData } from '../../../providers/conference-data';
import { UserData } from '../../../providers/user-data';
import { UserFilterPage } from '../user-filter/user-filter.page';
import { UsersService } from '../../../services/users.service';

@Component({
  selector: 'app-users-list',
  templateUrl: './users-list.page.html',
  styleUrls: ['./users-list.page.scss'],
})
export class UsersListPage implements OnInit {
  // Gets a reference to the list element
  @ViewChild('scheduleList', { static: true }) scheduleList: IonList;

  ios: boolean;
  dayIndex = 0;
  queryText = '';
  segment = 'all';
  excludeTracks: any = [];
  shownSessions: any = [];
  groups: any = [];
  confDate: string;
  showSearchbar: boolean;

  constructor(
    public alertCtrl: AlertController,
    public confData: ConferenceData,
    public loadingCtrl: LoadingController,
    public modalCtrl: ModalController,
    public router: Router,
    public routerOutlet: IonRouterOutlet,
    public toastCtrl: ToastController,
    public user: UserData,
    public config: Config,
    private userService: UsersService
  ) { }

  userslist = [
    {
      "id": 1,
      "Name": "Pins Prestilien",
      "Email": "pinsdev24@gmail.com",
      "Groups": [
        "Data Science"
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
      "id": 2,
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
      "id": 3,
      "Name": "Delano Roosvelt",
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
    },
    {
      "id": 4,
      "Name": "Youmie Yannick",
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
    },
    {
      "id": 5,
      "Name": "Ivan Ivan",
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

   key: string = 'userlist';
  
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
    this.updateSchedule();

    this.ios = this.config.get('mode') === 'ios';
  }

  tagclick () {
    
  }



  updateSchedule() {
    // Close any open sliding items when the schedule updates
    if (this.scheduleList) {
      this.scheduleList.closeSlidingItems();
    }

    this.confData.getTimeline(this.dayIndex, this.queryText, this.excludeTracks, this.segment).subscribe((data: any) => {
      this.shownSessions = data.shownSessions;
      this.groups = data.groups;
    });
  }

  async presentFilter() {
    const modal = await this.modalCtrl.create({
      component: UserFilterPage,
      swipeToClose: true,
      presentingElement: this.routerOutlet.nativeEl,
      componentProps: { excludedTracks: this.excludeTracks }
    });
    await modal.present();

    const { data } = await modal.onWillDismiss();
    if (data) {
      this.excludeTracks = data;
      this.updateSchedule();
    }
  }

  async addFavorite(slidingItem: HTMLIonItemSlidingElement, sessionData: any) {
    if (this.user.hasFavorite(sessionData.name)) {
      // Prompt to remove favorite
      this.removeFavorite(slidingItem, sessionData, 'Favorite already added');
    } else {
      // Add as a favorite
      this.user.addFavorite(sessionData.name);

      // Close the open item
      slidingItem.close();

      // Create a toast
      const toast = await this.toastCtrl.create({
        header: `${sessionData.name} was successfully added as a favorite.`,
        duration: 3000,
        buttons: [{
          text: 'Close',
          role: 'cancel'
        }]
      });

      // Present the toast at the bottom of the page
      await toast.present();
    }

  }

  async removeFavorite(slidingItem: HTMLIonItemSlidingElement, sessionData: any, title: string) {
    const alert = await this.alertCtrl.create({
      header: title,
      message: 'Would you like to remove this session from your favorites?',
      buttons: [
        {
          text: 'Cancel',
          handler: () => {
            // they clicked the cancel button, do not remove the session
            // close the sliding item and hide the option buttons
            slidingItem.close();
          }
        },
        {
          text: 'Remove',
          handler: () => {
            // they want to remove this session from their favorites
            this.user.removeFavorite(sessionData.name);
            this.updateSchedule();

            // close the sliding item and hide the option buttons
            slidingItem.close();
          }
        }
      ]
    });
    // now present the alert on top of all other content
    await alert.present();
  }

  async openSocial(network: string, fab: HTMLIonFabElement) {
    const loading = await this.loadingCtrl.create({
      message: `Posting to ${network}`,
      duration: (Math.random() * 1000) + 500
    });
    await loading.present();
    await loading.onWillDismiss();
    fab.close();
  }

  uploadCsv ( event ) {
    for (const file of event.target.files) {
      const fileData = new FormData();
      fileData.append("file", file);
      this.userService.uploadUsersFile( fileData )
        .subscribe( resp => {
          console.log( resp );            
        });
    }
}
}