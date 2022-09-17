import { Component, OnInit } from "@angular/core";
import { DomSanitizer } from "@angular/platform-browser";
import {
  Matching,
  MatchingGroupReq,
  MatchingReq,
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
  matchingRequest: MatchingReq | MatchingGroupReq;
  matchesSelected: Matching[] = [];

  constructor(
    private sanitizer: DomSanitizer,
    private matchingService: UsersService
  ) {}

  isrecurrence: boolean = false;
  selectedDays: String[] = [];
  selectedPattern: String;
  selectedEveryPattern: String;
  selectedWeek: String;
  allpatterns: any[] = [
    {
      value: "minutly",
      label:"Minutly"
    },
    {
      value: "hourly",
      label:"Hourly"
    },
    {
      value: "daily",
      label:"Daily"
    },
    {
      value: "weekly",
      label:"Weekly"
    },
    {
      value: "monthly",
      label:"Monthly"
    },
  ];

  everyPattern: any[] = [
    {
      value: "monday",
      label:"Monday"
    },
    {
      value: "tuesdy",
      label:"Tuesday"
    },
    {
      value: "wenesday",
      label:"Wenesday"
    },
    {
      value: "thursday",
      label:"Thursday"
    },
    {
      value: "friday",
      label:"Friday"
    },
    {
      value: "saturday",
      label:"Saturday"
    },
    {
      value: "sunday",
      label:"Sunday"
    },
  ];
  

  weeks: any[] = [
    {
      value: "1week",
      label:"1WEEK"
    },
    {
      value: "2week",
      label:"2WEEK"
    },
    {
      value: "3week",
      label:"3WEEK"
    },
    {
      value: "4week",
      label:"4WEEK"
    },
   
  ];
  selectedPeriodes: String[] = [];
  allPeriodes: String[] = [
    "ALL_DAY","ALL_MONTH","ALL_YEAR"
  ]
  oneMatchselected:Matching;

  noRecurentDate: String;
  startDate: String;
  endDate: String ="2023-05-17";
  time = "14:00";

  ngOnInit () {
    setTimeout( () => {
      this.noRecurentDate = new Date().toISOString();
      this.startDate = new Date().toISOString();
    } );
    
    this.matchings = history.state.matchings;

    this.matchings?.forEach((match) =>
      match.users.forEach((user) => {
        if (user.avatar) {
          user.avatar = this.sanitizer.bypassSecurityTrustHtml(
            user.avatar["changingThisBreaksApplicationSecurity"]
          );
        }
      })
    );
    this.matchingRequest = history.state.matchingRequest;
  }

  sendMail() {
    this.matchingService
      .sendEmail(this.matchings)
      .subscribe((res) => console.log(res));
  }
  sendMailByMatch(match) {
    this.oneMatchselected = match;
  }

  reloadSelectedMatches() {
    if (this.matchesSelected.length > 1) {
      const users: User[] = [];
      const userswithavatar: User[] = [];
      const forbiddenConnections: User[][] = [];
      for (const match of this.matchesSelected) {
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
      }

      const req: MatchingReq = {
        size: this.matchesSelected[0].users.length,
        users: users,
        forbiddenConnections: forbiddenConnections,
      };
      this.matchingService.makematch(req).subscribe((matchings: Matching[]) => {
        matchings.forEach((match) =>
          match.users.forEach((user) => {
            user.avatar = userswithavatar.find(
              (usr) => usr.id === user.id
            )?.avatar;
          })
        );

        this.matchings = this.matchings.concat(matchings);
      });
    }
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
