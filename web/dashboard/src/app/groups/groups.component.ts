import {Component, OnInit} from '@angular/core';
import {MatTableDataSource} from "@angular/material";
import {HttpClient} from "@angular/common/http";

@Component({
  selector: 'app-groups',
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.css']
})

export class GroupsComponent implements OnInit {
  public name = "";
  public emails = "";
  public groups = new MatTableDataSource([]);
  public group_columns = ["name", "emails", "delete"];

  constructor(private http: HttpClient) {

  }

  ngOnInit() {
    this.getGroups()
  }

  getGroups() {
    this.http.post("api/v1/", {
      "command": "get_groups",
      "data": {}
    }).subscribe((data: []) => {
      this.groups.data = data;
    })
  }

  addGroup() {
    console.log(this.emails)
    if (this.emails !== "" && this.name !== "") {
      let emails = this.emails.split("\n").join(" ").replace(/\s\s+/g, ' ').split(" ");
      this.http.post("api/v1/", {
        "command": "add_group",
        "data": {
          "name": this.name,
          "emails": Array.from(new Set(emails)),
        }
      }).subscribe(() => {
        this.name = "";
        this.emails = "";
        this.getGroups()
      })
    }
  }

  removeGroup(id: string) {
    this.http.post("api/v1/", {
      "command": "remove_group",
      "data": {
        "id": id,
      }
    }).subscribe(() => {
      this.getGroups()
    })
  }
}
