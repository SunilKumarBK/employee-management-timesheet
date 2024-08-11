import { Component, OnInit } from '@angular/core';
import { AuthService } from '../service/auth.service';

@Component({
  selector: 'app-timesheet',
  templateUrl: './timesheet.component.html',
  styleUrls: ['./timesheet.component.css']
})
export class TimesheetComponent implements OnInit {
  timesheet: any[] = [];
  empId:number=1;

  constructor(private authService: AuthService) {}

  ngOnInit() {
    this.authService.getTimesheet(this.empId).subscribe(data => {
      this.timesheet = data[0];
      console.log(this.timesheet);
      
    });
  }
}
