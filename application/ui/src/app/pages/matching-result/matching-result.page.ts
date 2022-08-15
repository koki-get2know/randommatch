import { Component, OnInit } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { Matching, MatchingReq, User, UsersService } from '../../services/users.service';

@Component({
  selector: 'app-matching-result',
  templateUrl: './matching-result.page.html',
  styleUrls: ['./matching-result.page.scss'],
})
export class MatchingResultPage implements OnInit {

  matchings: Matching[] = [];
  matchesSelected: Matching[] = [];

  constructor(private sanitizer: DomSanitizer, private matchingService: UsersService) { 
  }

  ngOnInit () {
    this.matchings = history.state.matchings;
    this.matchings?.forEach(match => match.users.forEach(user => {
      if (user.avatar) {
        user.avatar = this.sanitizer.bypassSecurityTrustHtml(user.avatar['changingThisBreaksApplicationSecurity']);
      }
    }));

  }
  
  sendMail() {
    this.matchingService.sendEmail(this.matchings).subscribe(res => console.log(res));
  }

  reloadSelectedMatches() {
    if(this.matchesSelected.length > 1) {
      const users: User[] = [];
      const userswithavatar: User[] = [];
      const forbiddenConnections: User[][] = [];
      for (const match of this.matchesSelected) {
        this.matchings.splice(this.matchings.findIndex(m => m.id === match.id), 1);
        const forbiddenConnection: User[] = [];
        for (const user of match.users) {
          users.push({id: user.id, name: user.name});
          userswithavatar.push(user);
          forbiddenConnection.push({id: user.id, name: user.name});
        }
        forbiddenConnections.push(forbiddenConnection);
      }
      
      const req: MatchingReq = {
        size: this.matchesSelected[0].users.length,
        users: users,
        forbiddenConnections: forbiddenConnections
      };
      this.matchingService.makematch(req).subscribe(( matchings: Matching[] ) => {
        matchings.forEach(match => match.users.forEach(user => {
          user.avatar = userswithavatar.find(usr => usr.id === user.id)?.avatar;
        }));
        
        this.matchings = this.matchings.concat(matchings);
      });
    }
  }

  selectTuple(event: PointerEvent, match: Matching) {  
    if ((event.target as HTMLInputElement).checked === false /* checkbox is checked */) {
      this.matchesSelected.push( match );
    } else {
      const index = this.matchesSelected.findIndex(m => match.id === m.id);
      this.matchesSelected.splice(index, 1);
    }
  }

}


