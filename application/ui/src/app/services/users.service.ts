import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class UsersService {

  constructor ( private http: HttpClient ) { }
  
  urlApi = "http://koki2.com:8011";

  async makematch(formData) {
    return await this.http.post( `${ this.urlApi }/matching`, formData, {} )
      .pipe( map( data => data ) );
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
