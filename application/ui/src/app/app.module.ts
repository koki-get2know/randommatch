import {
  HttpClient,
  HttpClientModule,
  HTTP_INTERCEPTORS,
} from "@angular/common/http";
import { NgModule } from "@angular/core";
import { BrowserModule } from "@angular/platform-browser";
import { InAppBrowser } from "@ionic-native/in-app-browser/ngx";
import { IonicModule } from "@ionic/angular";
import { IonicStorageModule } from "@ionic/storage";
import { TranslateLoader, TranslateModule } from "@ngx-translate/core";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { ServiceWorkerModule } from "@angular/service-worker";
import { environment } from "../environments/environment";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";

import {
  MsalGuard,
  MsalModule,
  MsalInterceptor,
  MsalRedirectComponent,
} from "@azure/msal-angular";
import {
  BrowserCacheLocation,
  InteractionType,
  PublicClientApplication,
} from "@azure/msal-browser";
import { IonicSelectableModule } from "ionic-selectable";
import { BearerInterceptor } from "./http-interceptors/bearer-interceptor.service";
import { appConstants } from "./constants";
import { TranslateHttpLoader } from "@ngx-translate/http-loader";

const isIE =
  window.navigator.userAgent.indexOf("MSIE ") > -1 ||
  window.navigator.userAgent.indexOf("Trident/") > -1;

export function HttpLoaderFactory(http: HttpClient) {
  return new TranslateHttpLoader(http, "./assets/i18n/", ".json");
}

@NgModule({
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    FormsModule,
    IonicSelectableModule,
    ReactiveFormsModule,
    IonicModule.forRoot(),
    IonicStorageModule.forRoot(),
    ServiceWorkerModule.register("ngsw-worker.js", {
      enabled: environment.production,
    }),
    TranslateModule.forRoot({
      defaultLanguage: "en",
    }),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useFactory: HttpLoaderFactory,
        deps: [HttpClient],
      },
    }),
    MsalModule.forRoot(
      new PublicClientApplication({
        auth: {
          clientId: environment.clientId, // Application (client) ID from the app registration
          authority: environment.authority, // The Azure cloud instance and the app's sign-in audience (tenant ID, common, organizations, or consumers)
          redirectUri: environment.redirectUri, // This is your redirect URI
          postLogoutRedirectUri: environment.redirectUri,
        },
        cache: {
          cacheLocation: BrowserCacheLocation.LocalStorage,
          storeAuthStateInCookie: isIE, // Set to true for Internet Explorer 11
        },
      }),
      {
        interactionType: InteractionType.Redirect, // MSAL Guard Configuration
        authRequest: {
          scopes: appConstants.scopes,
        },
      },
      {
        interactionType: InteractionType.Redirect, // MSAL Interceptor Configuration
        protectedResourceMap: new Map([
          [appConstants.microsoftGraph, appConstants.scopes],
        ]),
      }
    ),
  ],
  declarations: [AppComponent],
  providers: [
    InAppBrowser,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: BearerInterceptor,
      multi: true,
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: MsalInterceptor,
      multi: true,
    },
    MsalGuard,
  ],
  bootstrap: [AppComponent, MsalRedirectComponent],
})
export class AppModule {}
