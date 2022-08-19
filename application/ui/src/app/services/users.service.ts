import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { map, shareReplay, switchMap } from "rxjs/operators";
import { Observable } from "rxjs";
import { environment } from "../../environments/environment";
import { LoremIpsum } from "lorem-ipsum";
import { DomSanitizer, SafeHtml } from "@angular/platform-browser";
import avatar from "animal-avatar-generator";
import { MsalService } from "@azure/msal-angular";
import { appConstants } from "../constants";

export class User {
  id: string;
  name?: string;
  groups?: string[];
  //angular only param
  avatar?: SafeHtml;
  isChecked?: boolean = false;
}

interface UsersRes {
  data: User[];
}
export interface MatchingReq {
  size: number;
  users: User[];
  forbiddenConnections?: User[][];
}

export interface MatchingGroupReq {
  size: number;
  groups: User[][];
  forbiddenConnections?: User[][];
}
export interface Matching {
  id: number;
  users: User[];
}
interface MatchingResponse {
  data: Matching[];
}
interface JobResponse {
  status: string;
}

@Injectable({
  providedIn: "root",
})
export class UsersService {
  constructor(
    private http: HttpClient,
    private sanitizer: DomSanitizer,
    private authService: MsalService
  ) {}

  urlApi = environment.serverBaseUrl;

  makematch(matchingReq: MatchingReq): Observable<Matching[]> {
    return this.http
      .post<MatchingResponse>(`${this.urlApi}/matchings`, matchingReq)
      .pipe(map((res) => res.data));
  }

  makematchgroup(matchingReq: MatchingGroupReq): Observable<Matching[]> {
    return this.http
      .post<MatchingResponse>(`${this.urlApi}/group-matchings`, matchingReq)
      .pipe(map((res) => res.data));
  }

  uploadCsv(formData) {
    return this.http
      .post(`${this.urlApi}/matching`, formData, {})
      .pipe(map((data) => data));
  }

  uploadUsersFile(fileData: FormData) {
    return this.authService
      .acquireTokenSilent({
        scopes: appConstants.scopes,
        account: this.authService.instance.getAllAccounts()[0],
        forceRefresh: false,
      })
      .pipe(
        map((authResult) => {
          const roles: string[] | undefined = authResult.idTokenClaims["roles"];
          if (roles && roles.length > 0) {
            const prefix = "Org.";
            const orgs = roles
              .filter((x) => x.startsWith(prefix))
              .map((x) => x.slice(prefix.length));
            if (orgs.length > 0) {
              fileData.append("organization", orgs[0]);
            }
          }
          return fileData;
        }),
        switchMap((formdata) =>
          this.http
            .post(`${this.urlApi}/upload-users`, formdata, {
              reportProgress: true,
              observe: "response",
            })
            .pipe(map((data) => data))
        )
      );
  }

  availabilyofusers(checkurl): Observable<JobResponse> {
    return this.http
      .get<JobResponse>(`${this.urlApi}${checkurl}`)
      .pipe(shareReplay(1));
  }

  getUsersdata(): Observable<User[]> {
    return this.http.get<UsersRes>(`${this.urlApi}/users`).pipe(
      map((res) => {
        if (res.data) {
          for (const user of res.data) {
            user.avatar = this.generateAvatarSvg();
          }
        }
        return res.data;
      }),
      shareReplay(1)
    );
  }

  sendEmail(matches: Matching[]) {
    return this.http.post(`${this.urlApi}/email-matches`, { matches });
  }

  get() {
    return this.http
      .get<any>(`${this.urlApi}/matching/`)
      .pipe(map((data) => data));
  }

  // https://stackoverflow.com/questions/521295/seeding-the-random-number-generator-in-javascript/47593316#47593316
  private cyrb128(str: string) {
    let h1 = 1779033703,
      h2 = 3144134277,
      h3 = 1013904242,
      h4 = 2773480762;
    for (let i = 0; i < str.length; i++) {
      let k = str.charCodeAt(i);
      h1 = h2 ^ Math.imul(h1 ^ k, 597399067);
      h2 = h3 ^ Math.imul(h2 ^ k, 2869860233);
      h3 = h4 ^ Math.imul(h3 ^ k, 951274213);
      h4 = h1 ^ Math.imul(h4 ^ k, 2716044179);
    }
    h1 = Math.imul(h3 ^ (h1 >>> 18), 597399067);
    h2 = Math.imul(h4 ^ (h2 >>> 22), 2869860233);
    h3 = Math.imul(h1 ^ (h3 >>> 17), 951274213);
    h4 = Math.imul(h2 ^ (h4 >>> 19), 2716044179);
    return [
      (h1 ^ h2 ^ h3 ^ h4) >>> 0,
      (h2 ^ h1) >>> 0,
      (h3 ^ h1) >>> 0,
      (h4 ^ h1) >>> 0,
    ];
  }

  private mulberry32(a: number) {
    return () => {
      let t = (a += 0x6d2b79f5);
      t = Math.imul(t ^ (t >>> 15), t | 1);
      t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
      return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
    };
  }

  generateAvatarSvg(): SafeHtml {
    const lorem = new LoremIpsum();
    const seed = this.cyrb128(lorem.generateWords(2));
    const rand = this.mulberry32(seed[0]);
    return this.sanitizer.bypassSecurityTrustHtml(
      avatar(`${rand()}`, { size: 40 })
    );
  }

  compareconnection(forbconnec1: User[], forbconnec2: User[]) {
    return forbconnec1.filter((element) => {
      return !forbconnec2.some((elt2) => element.id === elt2.id);
    });
  }

  removeConnection(usersconnexionforbidden: User[][], index) {
    usersconnexionforbidden.splice(index, 1);
  }
}
