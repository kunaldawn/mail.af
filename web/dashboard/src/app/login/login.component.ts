import {Component, OnInit} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {Router} from "@angular/router";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  public signature = "";

  constructor(private http: HttpClient, private router: Router) {
  }

  ngOnInit() {
  }

  auth() {
    this.http.post("/auth/login/", {
      "signature": this.signature
    }, {
      "headers": new HttpHeaders().set("X-Skip-Interceptor", '')
    }).subscribe((token: string) => {
      localStorage.setItem("af_token", token);
      this.router.navigate(["jobs/"])
    })
  }
}
