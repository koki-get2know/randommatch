import { Component, OnInit } from "@angular/core";
import { FormBuilder, Validators, FormGroup } from "@angular/forms";
import {
  UsersService,
  MatchingReq,
  User,
  Matching,
} from "../../services/users.service";
import { NavController, ToastController } from "@ionic/angular";
import { NavigationExtras, Router } from "@angular/router";
import { ColorsTags } from "../../constants";
import { TranslateService } from "@ngx-translate/core";
import { finalize } from "rxjs/operators";
import SwiperCore, { Pagination, Swiper } from "swiper";

SwiperCore.use([Pagination]);

@Component({
  selector: "app-matching-simple",
  templateUrl: "./matching-simple.page.html",
  styleUrls: ["./matching-simple.page.scss"],
})
export class MatchingSimplePage implements OnInit {
  matchingForm: FormGroup;
  totalSelected = 0;
  users: User[] = [];
  ColorsTags = ColorsTags;
  usersSelected: User[] = [];
  usersToRestrictSelected: User[] = [];

  forbiddenConnections: User[][] = [];
  preferredConnections: User[][] = [];

  isIndeterminate: boolean;
  masterCheck: boolean;
  checkBoxList: any;
  isLoading = false;

  private slides: Swiper;

  constructor(
    private formBuilder: FormBuilder,
    private matchService: UsersService,
    public navCtrl: NavController,
    private router: Router,
    public toastController: ToastController,
    private translate: TranslateService
  ) {}

  ngOnInit() {
    this.matchService.getUsersdata().subscribe((users) => {
      this.users = users;
    });
    this.initForm();
  }

  initForm() {
    this.matchingForm = this.formBuilder.group({
      matchingSize: [
        "",
        Validators.compose([Validators.required, Validators.min(2)]),
      ],
    });
  }

  get form() {
    return this.matchingForm.controls;
  }

  setSwiperInstance(swiper: Swiper) {
    this.slides = swiper;
  }

  prevSlide() {
    this.slides.slidePrev();
  }

  nextSlide() {
    if (
      this.slides.activeIndex === 1 &&
      this.usersSelected.length < Number(this.form.matchingSize.value)
    ) {
      this.presentToast("SELECTION_SIZE_CONSISTENCY_INSTR", {
        value: this.form.matchingSize.value,
      });
    } else if (
      this.usersSelected.length === Number(this.form.matchingSize.value)
    ) {
      this.forbiddenConnections = [];
      this.preferredConnections = [];
      this.randommatch();
    } else {
      //remove connections that are not part of the selected users
      this.forbiddenConnections = this.forbiddenConnections.filter((users) => {
        for (const user of users) {
          if (!this.usersSelected.find((u) => u.id === user.id)) {
            return false;
          }
        }
        return true;
      });
      this.preferredConnections = this.preferredConnections.filter((users) => {
        for (const user of users) {
          if (!this.usersSelected.find((u) => u.id === user.id)) {
            return false;
          }
        }
        return true;
      });
      this.slides.slideNext();
    }
  }

  scroll(el: HTMLElement) {
    el.scrollIntoView({ behavior: "smooth" });
  }

  checkMaster() {
    this.usersSelected = [];
    setTimeout(() => {
      if (this.masterCheck) {
        this.users.forEach((user) => {
          user.isChecked = this.masterCheck;
          const copy = { ...user };
          copy.isChecked = false;
          this.usersSelected.push(copy);
        });
        this.totalSelected = this.usersSelected.length;
      } else {
        this.users.forEach((user) => {
          user.isChecked = this.masterCheck;
          this.onRemoveusersSelected(user.id);
        });
        this.totalSelected = 0;
      }
    });
  }

  private addCopyInarray(users: User[], user: User) {
    const copy = { ...user };
    copy.isChecked = false;
    users.push(copy);
  }
  private markUserAsSelected(user: User) {
    this.addCopyInarray(this.usersSelected, user);
  }

  private markUserAsRestricted(user: User) {
    this.addCopyInarray(this.usersToRestrictSelected, user);
  }

  private connectionAlreadyExists(
    connection: User[],
    connections: User[][]
  ): boolean {
    let i = 0;
    while (i < connections.length) {
      let element = connections[i];
      if (element.length === connection.length) {
        const diffUser = this.matchService.compareconnection(
          element,
          connection
        );
        if (diffUser.length === 0) {
          return true;
        }
      }
      i++;
    }
    return false;
  }

  private removeUserIdFromArray(collection: User[], id: string) {
    const index = collection.findIndex((d) => d.id === id);
    if (index >= 0) {
      collection.splice(index, 1);
    }
  }
  // when user is unchecked, it should be remove
  onRemoveusersSelected(id: string) {
    this.removeUserIdFromArray(this.usersSelected, id);
  }

  private removeUserRestriction(id: string) {
    this.removeUserIdFromArray(this.usersToRestrictSelected, id);
  }

  selectToRestrict(event: PointerEvent, user: User) {
    if ((event.target as HTMLInputElement).checked === false) {
      this.markUserAsRestricted(user);
    } else {
      this.removeUserRestriction(user.id);
    }
  }

  checkEvent(event: PointerEvent, user: User) {
    const totalItems = this.users.length;
    if ((event.target as HTMLInputElement).checked === false) {
      this.markUserAsSelected(user);
      this.totalSelected++;
    } else {
      this.onRemoveusersSelected(user.id);
      this.totalSelected--;
    }
    if (this.totalSelected > 0 && this.totalSelected < totalItems) {
      //If even one item is checked but not all
      this.isIndeterminate = true;
      this.masterCheck = false;
    } else if (this.totalSelected === totalItems) {
      //If all are checked
      this.masterCheck = true;
      this.isIndeterminate = false;
    } else {
      //If none is checked
      this.isIndeterminate = false;
      this.masterCheck = false;
    }
  }

  async presentToast(
    message: string,
    params?: Object,
    durationInMs: number = 2000
  ) {
    const translatedMessage: string = await this.translate
      .get(message, params)
      .toPromise();
    const toast = await this.toastController.create({
      message: translatedMessage,
      duration: durationInMs,
    });
    toast.present();
  }

  private addRestrictedConnection(connections: User[][]) {
    if (this.usersToRestrictSelected.length > 1) {
      if (
        !this.connectionAlreadyExists(
          this.usersToRestrictSelected,
          this.preferredConnections
        ) &&
        !this.connectionAlreadyExists(
          this.usersToRestrictSelected,
          this.forbiddenConnections
        )
      ) {
        connections.push(this.usersToRestrictSelected);
        this.usersSelected.forEach((user) => {
          user.isChecked = false;
        });
        this.usersToRestrictSelected = [];
        this.presentToast("CONNECTION_ADDED", {}, 1000);
      } else {
        this.presentToast("CONNECTION_ALREADY_EXISTS");
      }
    } else {
      this.presentToast("SELECT_MORE_THAN_ONE_USER");
    }
  }

  forbid() {
    this.addRestrictedConnection(this.forbiddenConnections);
  }

  prefer() {
    if (
      this.usersToRestrictSelected.length ===
      Number(this.form.matchingSize.value)
    ) {
      this.addRestrictedConnection(this.preferredConnections);
    } else {
      this.presentToast("MINIMUM_NUM_IN_PREFERRED_CONNECTION", {
        value: this.form.matchingSize.value,
      });
    }
  }

  removeForbiddenConnection(index) {
    this.matchService.removeConnection(this.forbiddenConnections, index);
  }

  removePreferredConnection(index: number) {
    this.matchService.removeConnection(this.preferredConnections, index);
  }

  randommatch() {
    if (Number(this.form.matchingSize.value) < 2) {
      this.presentToast("QUESTION_GROUP_MIN_INSTR");
      return;
    }
    if (this.usersSelected.length < Number(this.form.matchingSize.value)) {
      this.presentToast("SELECTION_INCONSISTENT_WITH_SIZE");
      return;
    }
    const users: User[] = [];
    const forbiddenConnections: User[][] = [];
    for (const selected of this.usersSelected) {
      users.push({
        id: selected.id,
        name: selected.name,
        avatar: selected.avatar,
      });
    }
    for (const connection of this.forbiddenConnections) {
      const newConnection = [];
      for (let item of connection) {
        newConnection.push({ id: item.id, name: item.name });
      }
      forbiddenConnections.push(newConnection);
    }
    const matchingRequest: MatchingReq = {
      size: Number(this.form.matchingSize.value),
      users,
      forbiddenConnections,
    };

    this.isLoading = true;
    this.matchService
      .makematch(matchingRequest)
      .pipe(finalize(() => (this.isLoading = false)))
      .subscribe((matchings: Matching[]) => {
        if (matchings !== null) {
          matchings.forEach((match) =>
            match.users.forEach((user) => {
              user.avatar = matchingRequest.users.find(
                (usr) => usr.id === user.id
              )?.avatar;
            })
          );
          this.matchingresult(matchings, matchingRequest);
        } else {
          this.presentToast("NO_MATCHING_GENERATED");
        }
      });
  }

  // matching result
  matchingresult(matchings: Matching[], matchingRequest: MatchingReq) {
    const navigationExtras: NavigationExtras = {
      state: {
        matchings,
        matchingRequest,
      },
    };
    this.router.navigate(["/matching-result"], navigationExtras);
  }
}
