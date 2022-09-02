import { Component, OnInit, ViewChild } from "@angular/core";
import { FormBuilder, Validators, FormGroup } from "@angular/forms";
import {
  UsersService,
  User,
  Matching,
  MatchingGroupReq,
} from "../../services/users.service";
import { IonAccordionGroup, ToastController } from "@ionic/angular";
import { IonicSelectableComponent } from "ionic-selectable";
import { ColorsTags } from "../../constants";
import SwiperCore, { Pagination, Swiper } from "swiper";
import { NavigationExtras, Router } from "@angular/router";
import { TranslateService } from "@ngx-translate/core";
import { finalize } from "rxjs/operators";

SwiperCore.use([Pagination]);

@Component({
  selector: "app-matching-group",
  templateUrl: "./matching-group.page.html",
  styleUrls: ["./matching-group.page.scss"],
})
export class MatchingGroupPage implements OnInit {
  matchingForm: FormGroup;
  ColorsTags = ColorsTags;

  private slides: Swiper;

  groups: User[][] = [];
  users: User[] = [];
  usersSelectedForGroup: User[] = [];
  activeGroupToEdit = -1;
  isLoading = false;
  noUsersToShow = false;

  @ViewChild("addUsersToGroup", { static: false })
  addUsersToGroup: IonicSelectableComponent;
  @ViewChild("groupsAccordion", { static: false })
  groupsAccordion: IonAccordionGroup;

  forbiddenConnections: User[][] = [];
  preferredConnections: User[][] = [];
  usersToRestrictSelected: User[] = [];

  constructor(
    private formBuilder: FormBuilder,
    private matchService: UsersService,
    private router: Router,
    private translate: TranslateService,
    private toastController: ToastController
  ) {}

  ngOnInit() {
    this.matchService.getUsersdata().subscribe((users) => {
      this.users = users;
      this.noUsersToShow = users.length === 0;
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
    let total = 0;
    this.groups.forEach((group) => (total += group.length));
    if (
      this.slides.activeIndex === 1 &&
      total < Number(this.form.matchingSize.value)
    ) {
      this.presentToast("SELECTION_SIZE_CONSISTENCY_INSTR", {
        value: this.form.matchingSize.value,
      });
    } else if (total === Number(this.form.matchingSize.value)) {
      this.forbiddenConnections = [];
      this.preferredConnections = [];
      this.randommatch();
    } else {
      //remove connections that are not part of the selected users
      this.forbiddenConnections = this.forbiddenConnections.filter((users) => {
        for (const user of users) {
          if (
            !this.groups.find((group) => group.some((u) => u.id === user.id))
          ) {
            return false;
          }
        }
        return true;
      });
      this.preferredConnections = this.preferredConnections.filter((users) => {
        for (const user of users) {
          if (
            !this.groups.find((group) => group.some((u) => u.id === user.id))
          ) {
            return false;
          }
        }
        return true;
      });
      this.slides.slideNext();
    }
  }

  addGroup() {
    this.activeGroupToEdit = -1;
    this.addUsersToGroup.open();
  }

  removeGroup(event: PointerEvent, index: number) {
    event.stopPropagation();
    this.groups.splice(index, 1);
  }

  removeUserFromGroup(event: PointerEvent, groupIndex: number, index: number) {
    event.stopPropagation();
    this.groups[groupIndex].splice(index, 1);
    if (this.groups[groupIndex].length === 0) {
      this.groups.splice(groupIndex, 1);
    }
  }

  searchByTags(event: { component: IonicSelectableComponent; text: string }) {
    event.component.startSearch();
    const text = event.text.trim().toLocaleLowerCase();
    if (text) {
      event.component.items = this.users.filter(
        (user) =>
          user.name.toLocaleLowerCase().includes(text) ||
          user.tags.some((tag) => tag.toLocaleLowerCase().includes(text))
      );
    } else {
      event.component.items = this.users;
    }
    event.component.endSearch();
  }

  addNewUsersToGroup() {
    this.addUsersToGroup.confirm();
    if (this.usersSelectedForGroup.length > 0) {
      if (this.activeGroupToEdit === -1) {
        this.groups = this.groups
          .map((group: User[]) => {
            return group.filter(
              (user) =>
                !this.usersSelectedForGroup.some((u) => u.id === user.id)
            );
          })
          .filter((group) => group.length > 0);
        this.groups.push([...this.usersSelectedForGroup]);
        this.groupsAccordion.value = (this.groups.length - 1).toString();
      } else {
        this.groups = this.groups
          .map((group: User[], index: number) => {
            if (index === this.activeGroupToEdit) {
              this.groupsAccordion.value = index.toString();
              return [...this.usersSelectedForGroup];
            }
            return group.filter(
              (user) =>
                !this.usersSelectedForGroup.some((u) => u.id === user.id)
            );
          })
          .filter((group) => group.length > 0);
      }
    }
    this.addUsersToGroup.clear();
    this.addUsersToGroup.close();
  }

  editGroup(event: PointerEvent, index: number) {
    event.stopPropagation();
    this.activeGroupToEdit = index;
    this.usersSelectedForGroup = [...this.groups[index]];
    this.addUsersToGroup.open();
  }

  toggleAll() {
    this.addUsersToGroup.toggleItems(
      this.addUsersToGroup.itemsToConfirm.length === 0
    );
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

  private addCopyInarray(users: User[], user: User) {
    const copy = { ...user };
    copy.isChecked = false;
    users.push(copy);
  }

  private markUserAsRestricted(user: User) {
    this.addCopyInarray(this.usersToRestrictSelected, user);
  }

  private removeUserIdFromArray(collection: User[], id: string) {
    const index = collection.findIndex((d) => d.id === id);
    if (index >= 0) {
      collection.splice(index, 1);
    }
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
        this.groups.forEach((group) =>
          group.forEach((user) => {
            user.isChecked = false;
          })
        );
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

  scroll(el: HTMLElement) {
    el.scrollIntoView({ behavior: "smooth" });
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
    let total = 0;
    this.groups.forEach((group) => (total += group.length));
    if (total < Number(this.form.matchingSize.value)) {
      this.presentToast("SELECTION_INCONSISTENT_WITH_SIZE");
      return;
    }

    const matchingRequest: MatchingGroupReq = {
      size: Number(this.form.matchingSize.value),
      groups: this.groups,
      forbiddenConnections: this.forbiddenConnections,
    };

    this.isLoading = true;
    this.matchService
      .makematchgroup(matchingRequest)
      .pipe(finalize(() => (this.isLoading = false)))
      .subscribe((matchings: Matching[]) => {
        if (matchings !== null) {
          matchings.forEach((match) =>
            match.users.forEach((user) => {
              for (const users of matchingRequest.groups) {
                for (const usr of users) {
                  if (user.id === usr.id) {
                    user.avatar = usr.avatar;
                    break;
                  }
                }
              }
            })
          );
          this.matchingresult(matchings, matchingRequest);
        } else {
          this.presentToast("NO_MATCHING_GENERATED");
        }
      });
  }

  // matching result
  matchingresult(
    matchings: Matching[],
    matchingGroupRequest: MatchingGroupReq
  ) {
    const navigationExtras: NavigationExtras = {
      state: {
        matchings,
        matchingGroupRequest,
      },
    };
    this.router.navigate(["/matching-result"], navigationExtras);
  }
}
