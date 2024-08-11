import { Component, inject,OnInit } from '@angular/core';
import { GoogleChartsModule, ChartType } from 'angular-google-charts';
import { CommonModule } from '@angular/common';
import { DepartmentService } from '../service/department.service';
import { EmployeeService } from '../service/employee.service';
import { MatDialog } from '@angular/material/dialog';
import { DialogComponent } from '../dialogBody/dialog/dialog.component';
import { MatIcon } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import { Location } from '@angular/common';

@Component({
  selector: 'app-organization-chart',
  standalone: true,
  imports: [GoogleChartsModule, CommonModule,MatIcon,RouterLink],
  templateUrl: './organization-chart.component.html',
  styleUrls: ['./organization-chart.component.css']
})
export class OrganizationChartComponent implements OnInit {
  chartData: any[];
  chartOptions: any;
  hierarchyData: any[];
  charttype = ChartType;
  selectedProfile: any;

  constructor(private service: DepartmentService, private matDialog: MatDialog,private location: Location) {
    this.chartData = [];  // Initialize as empty array
    this.chartOptions = {
      allowHtml: true,
      allowCollapse: true,
      node: {
        colors: ['#e0440e', '#e6693e', '#ec8f6e', '#f3b49f', '#f6c7b6'],
        size: 'medium',
      },
      tooltip: { isHtml: true },
      layout: 'horizontal',
      legend: { position: 'none' }
    };
    

  }
  goBack(): void {
    this.location.back();
  }

  ngOnInit(): void {
    this.service.getHierarchyData().subscribe({
      next: (data: any[]) => {
        this.hierarchyData = data.map(item => {
          if (item.tech_lead) {
            item.tech_lead = JSON.parse(item.tech_lead);
          }
          return item;
        });
        console.log(this.hierarchyData, 'Received hierarchy data');
        this.buildChartData();
      }
    });
  }
  
  



  private buildChartData(): void {
    // Initialize the chart with static nodes for CEO and Assistant CEO
    this.chartData = [
      [{ v: 'CEO', f: 'John Doe<div style="color:red; font-style:italic">CEO</div>' }, '', 'CEO'],
      [{ v: 'Assistant CEO', f: 'Jane Smith<div style="color:blue; font-style:italic">Assistant CEO</div>' }, 'CEO', 'Assistant CEO']
    ];
  
    // Create a map to store manager and department data
    const managerMap: { [key: string]: any } = {};
  
    // Process each employee in the hierarchy
    this.hierarchyData.forEach(employee => {
      const employeeName = `${employee.firstName} ${employee.lastName}`;
      const department = employee.department;
      const role = employee.roleName;
      const roleInfo = `<div style="color:gray">${role}</div>`;
      const managerName = employee.managerName || 'Assistant CEO';
  
      // If the employee is a manager and has no tech lead, they lead the department
      if (role === 'Manager' && employee.tech_lead === 0) {
        if (!managerMap[department]) {
          managerMap[department] = {
            manager: employeeName,
            employees: []
          };
  
          // Add manager to chart under Assistant CEO
          this.chartData.push([
            {
              v: employeeName,
              f: `${employeeName}<div style="color:green; font-style:italic">Manager</div><div style="color:blue;">${department}</div>`
            },
            'Assistant CEO',
            'Manager'
          ]);
        }
      }
  
      // Check if employee is under a tech lead and add them accordingly
      if (employee.tech_lead !== 0) {
        const manager = this.hierarchyData.find(emp => emp.employee_id === employee.tech_lead);
        if (manager) {
          const managerName = `${manager.firstName} ${manager.lastName}`;
          if (managerName) {
            managerMap[managerName].employees.push({
              name: employeeName,
              roleInfo: roleInfo
            });
  
            // Add employee under their manager
            this.chartData.push([
              { v: employeeName, f: `${employeeName}${roleInfo}` },
              managerName,
              'Employee'
            ]);
          }
        }
      }
    });
  
    console.log(this.chartData, 'Chart Data');
  }
  

  onSelect(event: any): void {
    const selection = event.selection;
    if (selection.length > 0) {
      const row = selection[0].row;

      // Ensure the row index is within bounds
      if (row >= 0 && row < this.chartData.length) {
        const node = this.chartData[row];
        const name = node[0].v;

        // Find the profile based on the selected name
        const profile = this.hierarchyData.find(emp => `${emp.firstName} ${emp.lastName}` === name);

        this.selectedProfile = profile || { firstName: 'Unknown', lastName: 'Unknown', roleName: 'Unknown', department: 'Unknown', managerName: 'Unknown' };
        console.log('Selected Profile:', this.selectedProfile);
        if(this.selectedProfile.empId!=undefined){
        this.openDialog(this.selectedProfile);
        }
      } else {
        console.log('Invalid row index selected');
      }
    }
  }

  passid: EmployeeService = inject(EmployeeService)

  openDialog(id: any) {
    console.log(id);

    this.passid.onPass(id.empId);
    this.matDialog.open(DialogComponent, {
      width: '400px',
    })
  }


}

