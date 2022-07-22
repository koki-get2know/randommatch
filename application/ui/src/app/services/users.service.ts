import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';

export class User {
  userId: string;
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
  
  urlApi = "http://localhost:8080";

  makematch(matchingReq: MatchingReq) : Observable<Matching[]> {
    return this.http.post<MatchingResponse>( `${ this.urlApi }/matchings`, matchingReq, {} )
      .pipe( map( res => res.data ) );
  }

  async uploadCsv(formData) {
    return await this.http.post( `${ this.urlApi }/matching`, formData, {} )
      .pipe( map( data => data ) );
  }

  async get() {
    return await this.http.get<any>(`${this.urlApi}/matching/`)
      .pipe(map(data => data));
  }

}
