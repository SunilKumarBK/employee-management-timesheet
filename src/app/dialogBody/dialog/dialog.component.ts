import { Component, inject, Input, input, OnInit } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogActions, MatDialogClose, MatDialogModule, MatDialogContent } from '@angular/material/dialog';
import { EmployeeService } from '../../service/employee.service';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { DepartmentService } from '../../service/department.service';
import { ActivatedRoute } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { error } from 'console';
import { MatTooltipModule } from '@angular/material/tooltip'; 

@Component({
  selector: 'app-dialog',
  standalone: true,
  imports: [MatButtonModule, MatDialogModule, MatDialogActions, MatDialogClose, MatDialogContent,RouterLink,CommonModule,MatIconModule,MatTooltipModule],
  templateUrl: './dialog.component.html',
  styleUrl: './dialog.component.css'
})
export class DialogComponent implements OnInit {
  // private _employeedata: any;
  getbyiddata: any;
  totalDuration:any;
  hierarchyData:any='';
  getdocuments:boolean;

  constructor(private service: DepartmentService,private _getid: EmployeeService,private route: ActivatedRoute, private router: Router) { }
  id: any;
  showAssignButton:boolean;
  showassigngetbyid:any;
  showassignhierarchy:any;
  showAssign:boolean;


  ngOnInit() {

   // Subscribe to the route's URL changes
   this.router.events.subscribe(() => {
    this.checkRoute();
  });
  
  // Initial check for the current route
  this.checkRoute();


// Method to check the current route and set the button visibility


    // this.testing();

    this.id = this._getid.idtodialog;
    console.log(this.id,'id')
    this._getid.getbyid(this.id).subscribe({
      next: (data: any) => {
        this.getbyiddata = data;
        this.showassigngetbyid =  this.getbyiddata.employee.empId;
        console.log(this.showassigngetbyid,'showassign')
        console.log(this.getbyiddata, 'getbydata'); // Log received data for debugging
        if (this.getbyiddata.company && this.getbyiddata.company.length > 0) {
          const durations = this.getbyiddata.company.map(company => company.duration);
          this.totalDuration = this.sumDurations(durations);
          console.log(this.totalDuration, 'totalduration');
        }
      },
      error: (error: any) => {
        console.error('There was an error!', error);
      },
      complete: () => {
        console.log('Request complete');
      }
    });

    //
   
      this.service.getHierarchyData().subscribe({
        next: (data: any[]) => {
          this.hierarchyData = data;
          console.log(this.hierarchyData, 'received hierarchy data');
          this.showassignhierarchy =  this.hierarchyData[this.id].empId;
          console.log(this.showassignhierarchy,'showassign')
 
        }
      });


      
if(this.showassigngetbyid==this.showassignhierarchy){
  this.showAssign=true;
  console.log(this.showAssign,'showassignn');
  
}else{
  this.showAssign=false;
  console.log(this.showAssign,'showassignn');

}

this._getid.getDocuments(this.id).subscribe({
next:(data:any)=>{
  console.log(data,'data');
  this.getdocuments=false;
},
error:(error:Error)=>{
  console.log(error);
  this.getdocuments=true;

  
}
})
     
  

  };

  private checkRoute(): void {
    const currentRoute = this.router.url;
    this.showAssignButton = currentRoute.includes('hierarchichal-chart');
    console.log(this.showAssignButton, 'showAssignButton');
  }



parseDuration(duration: string): { years: number, months?: number, days?: number } {
  const regex = /(?:(\d+)\s*years?)?\s*(?:(\d+)\s*months?)?\s*(?:(\d+)\s*days?)?/i;
  const match = duration.match(regex);

  if (!match) {
    console.log(`No match found for duration: ${duration}`);
    return { years: 0, months: 0, days: 0 };
  }

  const years = match[1] ? parseInt(match[1], 10) : 0;
  const months = match[2] ? parseInt(match[2], 10) : 0;
  const days = match[3] ? parseInt(match[3], 10) : 0;

  console.log(`Parsed duration: ${duration}`, { years, months, days });

  if(years && months && days !=0){
  return { years, months, days };
    
  }
  return { years };
}

  




 sumDurations(durations: string[]) {
  let totalYears = 0, totalMonths = 0, totalDays = 0;

  durations.forEach(duration => {
    const { years, months, days } = this.parseDuration(duration);
    totalYears += years;
    totalMonths += months;
    totalDays += days;
  });

  // Handle overflow
  while (totalDays >= 30) {
    totalDays -= 30;
    totalMonths += 1;
  }

  while (totalMonths >= 12) {
    totalMonths -= 12;
    totalYears += 1;
  }
if(totalYears&&totalMonths&&totalDays !=0){
  return `${totalYears} years ${totalMonths} months ${totalDays} days`;
}
return `${totalYears} years`
}
  


}



