import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { MsalService } from '@azure/msal-angular';
import { Observable } from 'rxjs';
import { map, shareReplay, switchMap } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { appConstants } from '../constants';

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

  constructor(private http: HttpClient,private authService: MsalService) { }

  urlApi = environment.serverBaseUrl;

  getMatchingStats(): Observable<MatchingStat[]>{
    return this.getOrganizations().pipe(
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

  getOrganizations(): Observable<string[]> {
    return this.authService
      .acquireTokenSilent({
        scopes: appConstants.scopes,
        account: this.authService.instance.getAllAccounts()[0],
        forceRefresh: false,
      })
      .pipe(
        map((authResult) => {
          let orgs: string[];
          const roles: string[] | undefined = authResult.idTokenClaims["roles"];
          if (roles && roles.length > 0) {
            const prefix = "Org.";
            orgs = roles
              .filter((x) => x.startsWith(prefix))
              .map((x) => x.slice(prefix.length));
          }
          return orgs;
        })
      );
  }
}
