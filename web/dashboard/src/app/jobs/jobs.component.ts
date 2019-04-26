import {Component, OnInit} from '@angular/core';
import {FormControl} from "@angular/forms";
import {HttpClient} from "@angular/common/http";
import {MatTableDataSource} from "@angular/material";

@Component({
  selector: 'app-jobs',
  templateUrl: './jobs.component.html',
  styleUrls: ['./jobs.component.css']
})
export class JobsComponent implements OnInit {
  public groupsList = [];
  public groups = new FormControl();
  public imageFile = new FormControl();
  public name = "";
  public subject = "";
  public min = new Date();
  public start = new Date();
  public completed = new MatTableDataSource([]);
  public completed_columns = ["name", "start", "end", "receivers", "success", "failures", "delete", "show"];
  public pending = new MatTableDataSource([]);
  public pending_columns = ["name", "start", "running", "receivers", "delete"];

  constructor(private http: HttpClient) {
  }

  ngOnInit() {
    this.getGroups();
    this.getCompletedJobs();
    this.getPendingJobs();
    setInterval(() => {
      this.getGroups();
      this.getCompletedJobs();
      this.getPendingJobs();
    }, 10000);
  }

  getGroups() {
    this.http.post("api/v1/", {
      "command": "get_groups",
      "data": {}
    }).subscribe((data: []) => {
      this.groupsList = data;
    })
  }

  getPendingJobs() {
    this.http.post("api/v1/", {
      "command": "get_pending_jobs",
      "data": {}
    }).subscribe((data: []) => {
      this.pending.data = data;
    })
  }

  getCompletedJobs() {
    this.http.post("api/v1/", {
      "command": "get_completed_jobs",
      "data": {}
    }).subscribe((data: []) => {
      this.completed.data = data;
    })
  }

  addJob() {
    if (this.name !== "" && this.subject !== "" && this.imageFile.value.files.length === 1 && this.groups.value && this.groups.value.length != 0) {
      this.getBase64(this.imageFile.value.files[0]).then((data) => {
        this.http.post("api/v1/", {
          "command": "create_job",
          "data": {
            "name": this.name,
            "groups": this.groups.value,
            "subject": this.subject,
            "image": data,
            "start": this.start.getTime() * 1000000
          }
        }).subscribe(() => {
          this.name = "";
          this.subject = "";
          this.groups.reset();
          this.imageFile.reset();
          this.getPendingJobs();
          this.getCompletedJobs();
        })
      });
    }
  }

  removeJob(id: string) {
    this.http.post("api/v1/", {
      "command": "remove_job",
      "data": {
        "id": id,
      }
    }).subscribe(() => {
      this.getPendingJobs();
      this.getCompletedJobs();
    })
  }

  getDate(time) {
    let ms = time / 1000000;
    let ds = new Date(ms);
    return ds.toLocaleString()
  }

  getBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = error => reject(error);
    });
  }

  openLogs(id: string) {
    window.open("logs/" + id, "_blank");
  }
}
