import {
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from "@angular/common/http";
import { Injectable } from "@angular/core";
import { MsalService } from "@azure/msal-angular";
import { Observable, of } from "rxjs";
import { map, switchMap } from "rxjs/operators";
import { appConstants } from "../constants";

@Injectable({
  providedIn: "root",
})
export class BearerInterceptor implements HttpInterceptor {
  constructor(private authService: MsalService) {}

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    if (req.url.startsWith(appConstants.microsoftGraph)) {
      return next.handle(req);
    }
    return this.getAuthorizationHeader().pipe(
      switchMap((idToken) => {
        const authReq = req.clone({
          setHeaders: {
            Authorization: `Bearer ${idToken}`,
          },
        });
        return next.handle(authReq);
      })
    );
  }

  private getAuthorizationHeader(): Observable<string> {
    return this.authService
      .acquireTokenSilent({
        scopes: appConstants.scopes,
        account: this.authService.instance.getAllAccounts()[0],
        forceRefresh: false,
      })
      .pipe(
        switchMap((authResult) => {
          if (
            new Date().getTime() - (authResult.idTokenClaims["exp"] as number) <
            10000
          ) {
            return of(authResult.idToken);
          } else {
            return this.authService
              .acquireTokenSilent({
                scopes: appConstants.scopes,
                account: this.authService.instance.getAllAccounts()[0],
                forceRefresh: true,
              })
              .pipe(map((authResult) => authResult.idToken));
          }
        })
      );
  }
}
