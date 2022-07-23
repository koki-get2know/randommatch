import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

export class User {
  userId: string;
  avatar?: string;
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

  constructor ( private http: HttpClient ) { }
  
  urlApi = environment.serverBaseUrl;

  makematch(matchingReq: MatchingReq) : Observable<Matching[]> {
    return this.http.post<MatchingResponse>( `${ this.urlApi }/matchings`, matchingReq)
      .pipe( map( res => res.data ) );
  }

  uploadCsv(formData) {
    return this.http.post( `${ this.urlApi }/matching`, formData, {} )
      .pipe( map( data => data ) );
  }

  get() {
    return this.http.get<any>(`${this.urlApi}/matching/`)
      .pipe(map(data => data));
  }

}
