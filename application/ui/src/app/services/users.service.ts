import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { LoremIpsum } from 'lorem-ipsum';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import avatar from 'animal-avatar-generator';

export class User {
  userId: string;
  avatar?: SafeHtml;
}
export interface MatchingReq {
  size: number;
  users: User[];
  forbiddenConnections?: User[][];
}

export interface Matching {
  id: number;
  users: User[];
}
interface MatchingResponse {
  data: Matching[];
}

@Injectable({
  providedIn: 'root'
})
export class UsersService {

  constructor ( private http: HttpClient, private sanitizer: DomSanitizer ) { }
  
  urlApi = environment.serverBaseUrl;

  makematch(matchingReq: MatchingReq) : Observable<Matching[]> {
    return this.http.post<MatchingResponse>( `${ this.urlApi }/matchings`, matchingReq)
      .pipe( map( res => res.data ) );
  }

  uploadUsersFile(fileData) {
    return this.http.post( `${ this.urlApi }/upload-users`, fileData, {
      reportProgress: true,
      observe: 'events'
    } )
      .pipe( map( data => data ) );
  }

  sendEmail(matches: Matching[]) {
    return this.http.post(`${ this.urlApi }/email-matches`, {matches,});
  }

  get() {
    return this.http.get<any>(`${this.urlApi}/matching/`)
      .pipe(map(data => data));
  }

  // https://stackoverflow.com/questions/521295/seeding-the-random-number-generator-in-javascript/47593316#47593316
  private cyrb128(str: string) {
    let h1 = 1779033703, h2 = 3144134277,
        h3 = 1013904242, h4 = 2773480762;
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
    return [(h1^h2^h3^h4)>>>0, (h2^h1)>>>0, (h3^h1)>>>0, (h4^h1)>>>0];
  }

  private mulberry32(a: number) {
    return () => {
      let t = a += 0x6D2B79F5;
      t = Math.imul(t ^ t >>> 15, t | 1);
      t ^= t + Math.imul(t ^ t >>> 7, t | 61);
      return ((t ^ t >>> 14) >>> 0) / 4294967296;
    }
  }

  generateAvatarSvg() : SafeHtml {
    const lorem = new LoremIpsum();
    const seed = this.cyrb128(lorem.generateWords(2));
    const rand = this.mulberry32(seed[0]);
    return this.sanitizer.bypassSecurityTrustHtml(avatar(`${rand()}`, { size: 40 }));
    
  }
}
