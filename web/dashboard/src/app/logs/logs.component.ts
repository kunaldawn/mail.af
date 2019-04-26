import {Component, OnInit} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {ActivatedRoute} from "@angular/router";
import {MatTableDataSource} from "@angular/material";

@Component({
  selector: 'app-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.css']
})

export class LogsComponent implements OnInit {
  public id = "";
  public logs = new MatTableDataSource([]);
  public log_columns = ["email", "time", "success"];

  constructor(private route: ActivatedRoute, private http: HttpClient) {
  }

  ngOnInit() {
    this.id = this.route.snapshot.paramMap.get('id');
    this.getLogs()
  }

  getLogs() {
    this.http.post("api/v1/", {
      "command": "get_logs",
      "data": {
        "id": this.id,
      }
    }).subscribe((data: []) => {
      this.logs.data = data
    })
  }

  getDate(time) {
    let ms = time / 1000000;
    let ds = new Date(ms);
    return ds.toLocaleString()
  }
}
