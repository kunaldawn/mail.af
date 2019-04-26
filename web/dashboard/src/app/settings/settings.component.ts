import {Component, OnInit} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {MatTableDataSource} from "@angular/material";

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.css']
})
export class SettingsComponent implements OnInit {
  public email = "";
  public password = "";
  public senders = new MatTableDataSource([]);
  public senders_columns = ["email", "delete"];


  constructor(private http: HttpClient) {

  }

  ngOnInit() {
    this.getSenders()
  }

  getSenders() {
    this.http.post("api/v1/", {
      "command": "get_senders",
      "data": {}
    }).subscribe((data: []) => {
      console.log(data);
      this.senders.data = data;
    })
  }

  addSender() {
    if (this.email !== "" && this.password !== "") {
      this.http.post("api/v1/", {
        "command": "add_sender",
        "data": {
          "email": this.email,
          "password": this.password,
        }
      }).subscribe(() => {
        this.email = "";
        this.password = "";
        this.getSenders()
      })
    }
  }

  removeSender(id: string) {
    this.http.post("api/v1/", {
      "command": "remove_sender",
      "data": {
        "id": id,
      }
    }).subscribe(() => {
      this.getSenders()
    })
  }
}
