import {Component, OnInit} from '@angular/core';
import {LoadingBarService} from "@ngx-loading-bar/core";
import {Router} from "@angular/router";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent implements OnInit {
  constructor(public loader: LoadingBarService, private router: Router) {
  }

  ngOnInit() {
  }

  openJobs() {
    this.router.navigate(['/jobs']);
  }

  openGroups() {
    this.router.navigate(['/groups']);
  }

  openSettings() {
    this.router.navigate(['/settings']);
  }

  logout() {
    localStorage.removeItem("af_token");
    this.router.navigate(["login/"])
  }
}
