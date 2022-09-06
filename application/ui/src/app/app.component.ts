import {
  Component,
  Inject,
  OnDestroy,
  OnInit,
  ViewEncapsulation,
} from "@angular/core";
import { Router } from "@angular/router";
import { SwUpdate } from "@angular/service-worker";

import { MenuController, Platform, ToastController } from "@ionic/angular";

import { StatusBar } from "@capacitor/status-bar";
import { SplashScreen } from "@capacitor/splash-screen";

import { Storage } from "@ionic/storage";

import {
  MsalBroadcastService,
  MsalGuardConfiguration,
  MsalService,
  MSAL_GUARD_CONFIG,
} from "@azure/msal-angular";
import {
  EventMessage,
  EventType,
  InteractionStatus,
  RedirectRequest,
} from "@azure/msal-browser";
import { filter, takeUntil } from "rxjs/operators";
import { Subject } from "rxjs";
import { environment } from "../environments/environment";
import { TranslateService } from "@ngx-translate/core";
@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.scss"],
  encapsulation: ViewEncapsulation.None,
})
export class AppComponent implements OnInit, OnDestroy {
  appPages = [
    {
      title: "USERS",
      url: "/users-list",
      icon: "people",
    },
    {
      title: "MATCHING",
      url: "matching",
      icon: "calendar",
    },
    {
      title: "STATISTICS",
      url: "statistics",
      icon: "bar-chart"
    }
  ];
  loggedIn = false;
  dark = false;
  isFrench = false;
  isIframe = false;
  private readonly _destroying$ = new Subject<void>();

  constructor(
    private menu: MenuController,
    private platform: Platform,
    private router: Router,
    private storage: Storage,
    private swUpdate: SwUpdate,
    private toastCtrl: ToastController,
    @Inject(MSAL_GUARD_CONFIG) private msalGuardConfig: MsalGuardConfiguration,
    private broadcastService: MsalBroadcastService,
    private authService: MsalService,
    private translate: TranslateService
  ) {
    this.initializeApp();
  }

  ngOnInit() {
    this.isIframe = window !== window.parent && !window.opener;
    this.broadcastService.msalSubject$
      .pipe(
        filter((msg: EventMessage) => msg.eventType === EventType.LOGIN_SUCCESS)
      )
      .subscribe((result: EventMessage) => {
        console.log(result);
      });

    this.broadcastService.inProgress$
      .pipe(
        filter(
          (status: InteractionStatus) => status === InteractionStatus.None
        ),
        takeUntil(this._destroying$)
      )
      .subscribe(() => {
        this.setLoggedIn();
      });

    this.swUpdate.versionUpdates.subscribe(async (res) => {
      const translation: { [key: string]: string } = await this.translate
        .get(["RELOAD", "UPDATE_AVAILABLE"])
        .toPromise();
      let message = "";
      let text = "";
      if (translation) {
        message = translation["UPDATE_AVAILABLE"] || "New version available";
        text = translation["RELOAD"] || "reload";
      }
      const toast = await this.toastCtrl.create({
        message: message,
        position: "bottom",
        buttons: [
          {
            role: "cancel",
            text: text,
          },
        ],
      });

      await toast.present();

      toast
        .onDidDismiss()
        .then(() => this.swUpdate.activateUpdate())
        .then(() => window.location.reload());
    });
  }

  initializeApp() {
    this.platform.ready().then(() => {
      if (this.platform.is("hybrid")) {
        StatusBar.hide();
        SplashScreen.hide();
      }
      this.translate.setDefaultLang("en");

      if (window.Intl && typeof window.Intl === "object") {
        const lang = navigator.language.substring(0, 2);
        this.isFrench = lang === "fr";
        this.translate.use(lang);
      } else {
        this.translate.use("en");
        this.isFrench = false;
      }
    });
  }

  changeLanguage() {
    if (this.isFrench) {
      this.translate.setDefaultLang("en");
      this.translate.use("fr");
    } else {
      this.translate.setDefaultLang("en");
      this.translate.use("en");
    }
  }
  login() {
    if (this.msalGuardConfig.authRequest) {
      this.authService.loginRedirect({
        ...this.msalGuardConfig.authRequest,
      } as RedirectRequest);
    } else {
      this.authService.loginRedirect();
    }
  }

  setLoggedIn() {
    this.loggedIn = this.authService.instance.getAllAccounts().length > 0;
  }

  logout() {
    this.authService.logoutRedirect({
      postLogoutRedirectUri: environment.redirectUri,
    });
  }

  openTutorial() {
    //this.menu.enable(false);
    this.storage.set("ion_did_tutorial", false);
    this.router.navigateByUrl("/tutorial");
  }

  ngOnDestroy() {
    this._destroying$.next(undefined);
    this._destroying$.complete();
  }
}
