import { HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { MsalService } from '@azure/msal-angular';
import { Observable } from 'rxjs';
import { map, switchMap } from 'rxjs/operators';
import { appConstants } from '../constants';

@Injectable( {
  providedIn: 'root'
} )
export class BearerInterceptor implements HttpInterceptor {

  constructor ( private authService: MsalService ) { }

  intercept ( req: HttpRequest<any>, next: HttpHandler ): Observable<HttpEvent<any>> {
    if ( req.url.startsWith( appConstants.microsoftGraph ) ) {
      return next.handle( req );
    }
    return this.getAuthorizationHeader()
      .pipe(
        switchMap( idToken => {
          const authReq = req.clone( {
            setHeaders: {
              Authorization: `Bearer ${ idToken }`,
            }
          } );
          return next.handle( authReq );
        } )
      );
  }

  private getAuthorizationHeader (): Observable<string> {
    return this.authService.acquireTokenSilent( {
      scopes: appConstants.scopes,
      account: this.authService.instance.getAllAccounts()[ 0 ],
      forceRefresh: false
    } ).pipe( map( authResult => authResult.idToken ) );
  }
}
