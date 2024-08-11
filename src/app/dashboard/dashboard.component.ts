import { Component, Injectable, OnInit } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { EmployeeService } from '../service/employee.service';
import { MatRippleModule } from '@angular/material/core';
import { DepartmentService } from '../service/department.service';
import { CommonModule, DatePipe } from '@angular/common';
// import { Console } from 'console';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [RouterLink, RouterLinkActive, MatCardModule, MatButtonModule, MatIconModule, MatRippleModule, CommonModule],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css',
  // providers: [DatePipe]
})

@Injectable({
  providedIn: 'root'
})
export class DashboardComponent implements OnInit {
  constructor(private service: EmployeeService, private departmentService: DepartmentService) { }
  employeeNo: any;
  addedmessage: string = 'default';
  deletedmessage: any;
  delmsg: any;
  temp: any;
  hierarchyData: any;
  departmentCounts: { [key: string]: number } = {};
  date:Date;

  negindex: number;
  delindex: number;


  ngOnInit(): void {
    // this.date = new Date();
    // this.employeeNo=0;
    this.reloadData();

    this.departmentService.getHierarchyData().subscribe({
      next: (data: any[]) => {
        this.hierarchyData = data;
        console.log(this.hierarchyData, 'received hierarchy data');

        if(this.hierarchyData!=null){
        //for loop
        for (let index = 0; index < this.hierarchyData.length; index++) {
          const department = this.hierarchyData[index].department;
          // Increment the count for the employee's department
          if (this.departmentCounts[department]) {
            this.departmentCounts[department]++;
          } else {
            this.departmentCounts[department] = 1;
          }
          console.log(this.departmentCounts);
        }
        console.log(this.departmentCounts['Sales']);
      }
    }


    });

  }

  // getShortDateTime(date: Date): string {
  //   return this.datePipe.transform(date, 'EEE, h:mm a') || '';
  // }

  reloadData() {
    this.service.getUsersList().subscribe(
      {
        next: (data: any[]) => {
          this.employeeNo = data;
          this.negindex = this.employeeNo.length - 1;
        },
        error: (error) => {
          console.error('There was an error!', error);
        },
        complete: () => {
          console.log('Request complete',this.employeeNo.length);
          this.addedmessage = this.service.addedmessage;
          this.date=new Date();

          this.deletedmessage = this.service.deletedmessage;
          this.delmsg = this.service.dellmsg;
          setTimeout(() => {
            this.addedmessage = "";
            this.deletedmessage = "";
            this.delmsg = "";
          }, 5000)
        }
      });

  }





}
