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
import SwiperCore, { Pagination, Swiper } from "swiper";

SwiperCore.use([Pagination]);

@Component({
  selector: "app-matching",
  templateUrl: "./matching.page.html",
  styleUrls: ["./matching.page.scss"],
})
export class MatchingPage implements OnInit {
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

  private slides: Swiper;

  constructor(
    private formBuilder: FormBuilder,
    private matchService: UsersService,
    public navCtrl: NavController,
    private router: Router,
    public toastController: ToastController
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

  nextSlide() {
    this.slides.allowSlideNext = true;
    this.slides.slideNext();
    this.slides.allowSlideNext = false;
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
  // when user is unchecked, it should be remove
  onRemoveusersSelected(id: string) {
    const index = this.usersSelected.findIndex((d) => d.id === id); //find index in your array
    const indexInUsers = this.users.findIndex((d) => d.id === id);
    if (indexInUsers >= 0) {
      this.users[indexInUsers].isChecked = false;
    }
    this.usersSelected.splice(index, 1);
  }

  removeUserRestriction(id: string) {
    const index = this.usersToRestrictSelected.findIndex((d) => d.id === id);
    if (index >= 0) {
      this.usersToRestrictSelected.splice(index, 1);
    }
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
    console.log(this.totalSelected);
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

  async presentToast(message) {
    const toast = await this.toastController.create({
      message: message,
      duration: 2000,
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
      } else {
        this.presentToast("this connection already exists!");
      }
    } else {
      this.presentToast("Please select more than one user!");
    }
  }

  forbid() {
    this.addRestrictedConnection(this.forbiddenConnections);
  }

  prefer() {
    this.addRestrictedConnection(this.preferredConnections);
  }

  removeForbiddenConnection(index) {
    this.matchService.removeConnection(this.forbiddenConnections, index);
  }

  removePreferredConnection(index: number) {
    this.matchService.removeConnection(this.preferredConnections, index);
  }

  randommatch() {
    if (Number(this.form.matchingSize.value) < 2) {
      this.presentToast("Matching size should be at least 2");
      return;
    }
    if (this.usersSelected.length < Number(this.form.matchingSize.value)) {
      this.presentToast("Users selected not consistent with matching size");
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

    this.matchService
      .makematch(matchingRequest)
      .subscribe((matchings: Matching[]) => {
        if (matchings !== null) {
          console.log(matchings);
          matchings.forEach((match) =>
            match.users.forEach((user) => {
              user.avatar = matchingRequest.users.find(
                (usr) => usr.id === user.id
              )?.avatar;
            })
          );
          this.matchingresult(matchings);
        } else {
          this.presentToast("No matchings generated!");
        }
      });
  }

  // matching result
  matchingresult(matchings: Matching[]) {
    const navigationExtras: NavigationExtras = {
      state: {
        matchings,
      },
    };
    this.router.navigate(["/matching-result"], navigationExtras);
  }
}
