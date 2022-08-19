import { Component, OnInit } from "@angular/core";
import { Router } from "@angular/router";
import { ToastController } from "@ionic/angular";

import { User, UsersService } from "../../../services/users.service";
import { ColorsTags } from "../../../constants";

@Component({
  selector: "app-users-list",
  templateUrl: "./users-list.page.html",
  styleUrls: ["./users-list.page.scss"],
})
export class UsersListPage implements OnInit {
  userslist: User[] = [];
  isloading: boolean = false;
  checkResponseUrl = "";
  ColorsTags = ColorsTags;
  constructor(
    public router: Router,
    public toastCtrl: ToastController,
    private userService: UsersService
  ) {}

  ngOnInit() {
    this.getuserList();
  }

  tagclick(event: Event) {
    event.stopPropagation();
  }

  uploadCsv(event: Event) {
    this.isloading = true;
    for (const file of event.target["files"]) {
      const fileData = new FormData();
      fileData.append("file", file);
      this.userService.uploadUsersFile(fileData).subscribe(
        (resp) => {
          if (resp.status === 202) {
            this.checkResponseUrl = resp.headers.get("Location");
            this.checkJobStatus();
          }
        },
        (_) => {
          this.isloading = false;
        }
      );
    }
  }

  checkJobStatus() {
    let responsestatus = "";
    const limitedInterval = setInterval(() => {
      this.userService
        .availabilyofusers(this.checkResponseUrl)
        .subscribe((resp) => {
          responsestatus = resp.status;
          if (responsestatus === "Done") {
            this.getuserList();
            this.isloading = false;
            clearInterval(limitedInterval);
          } else if (
            responsestatus !== "" &&
            responsestatus !== "Running" &&
            responsestatus !== "Pending"
          ) {
            this.isloading = false;
            clearInterval(limitedInterval);
            console.log(`responsestatus ${responsestatus} interval cleared!`);
          }
        });
    }, 500);
  }

  getuserList() {
    this.userService.getUsersdata().subscribe((users) => {
      this.userslist = users;
    });
  }
}
