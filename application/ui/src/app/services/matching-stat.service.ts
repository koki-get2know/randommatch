import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { MsalService } from '@azure/msal-angular';
import { Observable } from 'rxjs';
import { map, shareReplay, switchMap } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { UsersService } from './users.service';

export class MatchingStat {
  Id: string;
  NumGroups: number;
  NumConversations: number;
  NumFailed: number;
  CreatedAt: Date;
}

interface MatchingStatResponse {
  data: MatchingStat[];
}

@Injectable({
  providedIn: 'root'
})
export class MatchingStatService {

  constructor(private http: HttpClient,private authService: MsalService, private userService: UsersService) { }

  urlApi = environment.serverBaseUrl;

  getMatchingStats(): Observable<MatchingStat[]>{
    return this.userService.getOrganizations().pipe(
      map((orgs) => {
        let orga = "";
        if (orgs && orgs.length > 0) {
          orga = orgs[0];
        }
        return orga;
      }),
      switchMap((orga) =>
        this.http
          .get<MatchingStatResponse>(`${this.urlApi}/matchings-stats?organization=${orga}`)
          .pipe(
            map((res) => {
              return res.data;
            })
          )
      ),
      shareReplay(1)
    );
  }
}
