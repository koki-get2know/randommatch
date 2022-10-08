import { Component, OnInit } from "@angular/core";
import { ToastController } from "@ionic/angular";
import { TranslateService } from "@ngx-translate/core";
import { parseISO } from 'date-fns';

import {
  Matching,
  MatchingGroupReq,
  MatchingReq,
  ScheduleParam,
  SchedulingReq,
  User,
  UsersService,
} from "../../services/users.service";

@Component({
  selector: "app-matching-result",
  templateUrl: "./matching-result.page.html",
  styleUrls: ["./matching-result.page.scss"],
})
export class MatchingResultPage implements OnInit {
  matchings: Matching[] = [];
  matchingRequest: MatchingReq;
  matchingGroupRequest: MatchingGroupReq;
  matchesSelected: Matching[] = [];

  constructor(
    private matchingService: UsersService,
    private toastController: ToastController,
    private translate: TranslateService) {}

  isrecurrence: boolean = false;
  endedBySystem: boolean = false;
  duration: number = 15;
  selectedDays: string[] = [];
  selectedPattern: string;
  selectedEveryPattern: string[] = [];
  selectedWeek: string;
  allpatterns: any[] = [
    {
      value: "every_minute",
      label: "EVERY_MINUTE",
    },
    {
      value: "hourly",
      label: "HOURLY",
    },
    {
      value: "daily",
      label: "DAILY",
    },
    {
      value: "weekly",
      label: "WEEKLY",
    },
    {
      value: "monthly",
      label: "MONTHLY",
    },
  ];

  everyPattern: any[] = [
    {
      value: "Monday",
      label: "MONDAY",
    },
    {
      value: "Tuesday",
      label: "TUESDAY",
    },
    {
      value: "Wednesday",
      label: "WEDNESDAY",
    },
    {
      value: "Thursday",
      label: "THURSDAY",
    },
    {
      value: "Friday",
      label: "FRIDAY",
    },
    {
      value: "Saturday",
      label: "SATURDAY",
    },
    {
      value: "Sunday",
      label: "SUNDAY",
    },
  ];

  weeks: any[] = [
    {
      value: "first",
      label: "WEEK1",
    },
    {
      value: "second",
      label: "WEEK2",
    },
    {
      value: "third",
      label: "WEEK3",
    },
    {
      value: "fourth",
      label: "WEEK4",
    },
    {
      value: "last",
      label: "LAST_WEEK",
    },
  ];
  selectedPeriodes: String[] = [];
  allPeriodes: String[] = ["EVERY_DAY", "EVERY_MONTH", "EVERY_YEAR"];
  oneMatchselected: Matching;

  inviteDate: string;
  dateofDay: string;
  startDate: string;
  endDate: string;
  time: string;

  ngOnInit() {
    setTimeout(() => {
      this.inviteDate = new Date().toISOString();
      this.dateofDay = new Date().toISOString();
      this.startDate = new Date().toISOString();
      this.time = new Date().toISOString();
    } );
    
    

    this.matchings = history.state.matchings;
    this.matchingRequest = history.state.matchingRequest;
    this.matchingGroupRequest = history.state.matchingGroupRequest;
  }

  scheduleInvites() {
    if (this.isrecurrence) {
      this.makeSheduling();
    } else {
      this.sendMail();
    }
  }

  sendMail() {
    this.matchingService.sendEmail(this.matchings, this.inviteDate,this.duration).subscribe(
      (_: Object) => {
        this.presentToast("SCHEDULE_TAKEN_INTO_ACCOUNT");
      }
    );
  }

  makeSheduling () {
    let matchSize: number;
    let matchType: string;
    if ( this.matchingRequest !== undefined ) {
      matchSize = this.matchingRequest.size;
      matchType = "simple";
    }
    else if (this.matchingGroupRequest !== undefined) {
      matchSize = this.matchingGroupRequest.size;
      matchType = "group";
    }
    const sheduleParam: ScheduleParam = {
      time: parseISO(this.time).getUTCHours()+":"+parseISO(this.time).getUTCMinutes(),
      duration: this.duration+"",
      startDate: parseISO(this.startDate).toISOString(),
      endDate: parseISO(this.endDate).toISOString(),
      size: matchSize,
      matchingType: matchType,
      frequency: this.selectedPattern,
      week: this.selectedWeek,
      days: this.selectedEveryPattern
    };

    const schedulingRequest: SchedulingReq = {
      schedule: sheduleParam,
      ...( matchType === "simple" && { users: this.matchingRequest.users } ),
      ...( matchType === "group" && { group: this.matchingGroupRequest.groups } ),

    };
   
    this.matchingService.makesheduling( schedulingRequest ).subscribe(
      ( _: any ) => {
        this.presentToast("SCHEDULE_TAKEN_INTO_ACCOUNT");
      }
    );
  }


  groupReloadSelectedMatches() {
    if (this.matchesSelected.length > 1) {
      const groups: User[][] = [];
      const userswithavatar: User[] = [];
      const forbiddenConnections: User[][] = [];
      let index = 0;
      let position = 0;
      let map = new Map<number, User[]>();

      for (const match of this.matchesSelected) {
        if (index === 0) {
          position = this.matchings.findIndex((m) => m.id === match.id);
        }
        this.matchings.splice(
          this.matchings.findIndex((m) => m.id === match.id),
          1
        );
        const forbiddenConnection: User[] = [];
        for (const user of match.users) {
          let idx = this.matchingGroupRequest.groups.findIndex((group) =>
            group.some((u) => u.id === user.id)
          );
          if (map.has(idx)) {
            map.get(idx).push(user);
          } else {
            map.set(idx, [user]);
          }

          userswithavatar.push(user);
          forbiddenConnection.push({ id: user.id, name: user.name });
        }

        forbiddenConnections.push(forbiddenConnection);
        if (this.matchingGroupRequest.forbiddenConnections) {
          forbiddenConnections.push(
            ...this.matchingGroupRequest.forbiddenConnections
          );
        }
      }

      for (const [_, group] of map) {
        groups.push(group);
      }
      const req: MatchingGroupReq = {
        size: this.matchesSelected[0].users.length,
        groups,
        forbiddenConnections: forbiddenConnections,
      };
      this.matchingService
        .makematchgroup(req)
        .subscribe((matchings: Matching[]) => {
          if (matchings) {
            matchings.forEach((match) =>
              match.users.forEach((user) => {
                user.avatar = userswithavatar.find(
                  (usr) => usr.id === user.id
                )?.avatar;
              })
            );
            this.matchings.splice(position, 0, ...matchings);
          } else {
            this.presentToast("MATCH_POSSIBILITY_EXHAUSTED");
          }
        });
    }
  }
  sendMailByMatch(match) {
    this.oneMatchselected = match;
  }


  reloadSelectedMatches() {
    if (this.matchesSelected.length === 0) {
      this.matchesSelected = [...this.matchings];
    }
    if (this.matchingRequest) {
      this.reloadSimpleSelectedMatches();
      this.matchesSelected = [];
    } else {
      this.groupReloadSelectedMatches();
      this.matchesSelected = [];
    }
  }

  reloadSimpleSelectedMatches() {
    if (this.matchesSelected.length > 1) {
      const users: User[] = [];
      const userswithavatar: User[] = [];
      const forbiddenConnections: User[][] = [];
      let index = 0;
      let position = 0;
      for (const match of this.matchesSelected) {
        if (index === 0) {
          position = this.matchings.findIndex((m) => m.id === match.id);
        }
        this.matchings.splice(
          this.matchings.findIndex((m) => m.id === match.id),
          1
        );
        const forbiddenConnection: User[] = [];
        for (const user of match.users) {
          users.push({ id: user.id, name: user.name });
          userswithavatar.push(user);
          forbiddenConnection.push({ id: user.id, name: user.name });
        }
        forbiddenConnections.push(forbiddenConnection);
        if (this.matchingRequest.forbiddenConnections) {
          forbiddenConnections.push(
            ...this.matchingRequest.forbiddenConnections
          );
        }
        index++;
      }

      const req: MatchingReq = {
        size: this.matchesSelected[0].users.length,
        users: users,
        forbiddenConnections: forbiddenConnections,
      };
      this.matchingService.makematch(req).subscribe((matchings: Matching[]) => {
        if (matchings) {
          matchings.forEach((match) =>
            match.users.forEach((user) => {
              user.avatar = userswithavatar.find(
                (usr) => usr.id === user.id
              )?.avatar;
            })
          );
          this.matchings.splice(position, 0, ...matchings);
        } else {
          this.presentToast("MATCH_POSSIBILITY_EXHAUSTED");
        }
      });
    }
  }

  async presentToast(
    message: string,
    params?: Object,
    durationInMs: number = 15000
  ) {
    const translatedMessage: string = await this.translate
      .get(message, params)
      .toPromise();
    const toast = await this.toastController.create({
      message: translatedMessage,
      duration: durationInMs,
    });
    toast.present();
  }

  selectTuple(event: PointerEvent, match: Matching) {
    if (
      (event.target as HTMLInputElement).checked ===
      false /* checkbox is checked */
    ) {
      this.matchesSelected.push(match);
    } else {
      const index = this.matchesSelected.findIndex((m) => match.id === m.id);
      this.matchesSelected.splice(index, 1);
    }
  
  
  }
}
