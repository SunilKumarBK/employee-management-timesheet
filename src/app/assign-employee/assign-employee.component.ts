import { CommonModule } from '@angular/common';
import { Component, OnInit, inject } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { EmployeeService } from '../service/employee.service';
import { DepartmentService } from '../service/department.service';
import { MatButtonModule } from '@angular/material/button';
import { ToastrService } from 'ngx-toastr';
import { ActivatedRoute, Router } from '@angular/router';
import { error } from 'node:console';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-assign-employee',
  standalone: true,
  imports: [MatCardModule, CommonModule, FormsModule, ReactiveFormsModule, MatButtonModule, MatIconModule],
  templateUrl: './assign-employee.component.html',
  styleUrl: './assign-employee.component.css'
})
export class AssignEmployeeComponent implements OnInit {

  constructor(private _employeedata: EmployeeService, private dataService: DepartmentService, private router: Router, private route: ActivatedRoute) { }


  assigningForm: FormGroup;
  employeedata: any;
  filteredemployee:any;
  employeeid: any = '';
  departments: any[] = [];
  roles: any[] = [];
  managerName: string = '';
  managerId: any;
  assignData: any;
  toaster = inject(ToastrService);
  employeeId: any;
  selected_dept: string;
  selectedDept: string;
  filteredassemp: any;
  disableemployeeselect: boolean = false;
  hierarchyData: any;
  roledata:any;
  hierarchytData:any;
  empId:any;
filteredRole:any;

  managers: any[] = [];

  availableJobRoles: string[] = [];

  reloadData() {
    this._employeedata.getUsersList().subscribe(
      {
        next: (data: any[]) => {
          this.employeedata = data;
          console.log(this.employeedata, 'employeeassigndata');


          this.dataService.getHierarchyData().subscribe({
            next: (hierarchyData: any) => {
              this.hierarchyData = hierarchyData;
              console.log(this.hierarchyData, 'hierarchy data');

              const assignedEmployeeIds = this.hierarchyData?.map(item => item.empId);

// Filter out employees who are not assigned
const filteredDatatest = this.employeedata.filter(employee => !assignedEmployeeIds?.includes(employee.empId));
              this.filteredemployee=filteredDatatest

console.log(filteredDatatest, 'filtered unassigned employees');

      
              // // Filter employee data based on hierarchy data
              // const filteredData = this.employeedata.filter(employee => employee.empId !== this.hierarchyData[0].empId);
              // this.filteredemployee=filteredData
              // console.log(filteredData, 'filtered data');
            },
            error: (err) => {
              console.error('Error fetching hierarchy data:', err);
            }
          });
        },
        error: (err) => {
          console.error('Error fetching employee data:', err);
        }
      });


  }

  ngOnInit(): void {

    this.reloadData();


    this.assigningForm = new FormGroup({
      empId: new FormControl(''),
      employee_id: new FormControl('Select Employee', [Validators.required, Validators.minLength(4)]),
      dept_id: new FormControl('Select Department', [Validators.required]),
      role_id: new FormControl('Select Job Role', [Validators.required]),
      // manager: new FormControl({value:'',disabled:false}, [Validators.required]),
      tech_lead: new FormControl({ value: '', disabled: true }, [Validators.required])

    })
    this.assigningForm.get('employee_id')?.valueChanges.subscribe(empId => {
      this.assigningForm.get('empId')?.setValue(empId);
      this.empId=empId
    });

    this.dataService.getDepartments().subscribe(data => {
      this.departments = data;
      console.log(this.departments, 'departments');

    });

    this.assigningForm.get('dept_id')?.valueChanges.subscribe(deptId => {
      console.log(deptId, 'deptid');




      if (deptId) {
        this.dataService.getRolesByDepartment(deptId).subscribe(data => {
          this.roles = data;
          console.log(this.roles,'roless');
          this.assigningForm.get('role_id')?.valueChanges.subscribe(role_id => {
          this.dataService.getHierarchyDatabyid(role_id).subscribe({next:(data:any)=>
          {
            this.hierarchytData=data
            console.log(data)
            // if(!this.hierarchyData){
           
            // }else{
            // const secfilter = this.roles.filter(data =>data.role_id == this.hierarchytData.role_id)
            // this.filteredRole=secfilter
            // console.log(secfilter);
            // }
            
            
          }
        
        }
        )
      })
          
         
          
          // Reset job role if necessary
          this.assigningForm.patchValue({ role_id: 'Select Job Role' });
          // Fetch manager based on department ID
          // this.getManagerByDepartment(deptId);
          this.getEmployeeAsManager(deptId)
          this.dataService.getDepartmentsById(deptId).subscribe(data => {
            this.selected_dept = data[0].department;
            this.selectedDept = String(this.selected_dept)
            console.log(this.selectedDept, 'datadept');

          })

          // this.dataService.getDepartments().subscribe
          // this.departments(deptId)

          this.dataService.getHierarchyData().subscribe(data => {
            console.log(data, 'data');
            
            // Check if data is null or undefined
            if (!data) {
              console.error('No hierarchy data received');
              this.filteredRole = this.roles; // Handle this case as needed

              return;
            }
          
            this.hierarchyData = data;
            
          
            // Ensure hierarchyData is an array
            if (!Array.isArray(this.hierarchyData) || this.hierarchyData.length === 0) {
              console.error('Hierarchy data is not an array or is empty');
              this.filteredRole = this.roles; // Handle this case as needed
              return;
            }
          
            // Create a set of assigned role IDs from hierarchyData
            const assignedRoleIds = new Set(this.hierarchyData.map(item => item.role_id));
          
            // Filter out roles that have already been assigned
            const filteredRole = this.roles.filter(role => !assignedRoleIds.has(role.role_id));
            this.filteredRole = filteredRole;
            this.filteredRole=filteredRole
          
            console.log(filteredRole, 'filteredRole');
              // Assuming data is an array of objects with a 'department' property
            const filteredData = data.filter(item => item.department == this.selectedDept); // Change 'Sales' to the department you want to filter by
            this.filteredassemp = filteredData
            console.log(filteredData, 'Filtered Data');
          });


        });
      } else {
        this.assigningForm.patchValue({ role_id: 'Select Job Role' });
      }
      this._employeedata.setManager(deptId).subscribe({
        next: (data: any) => {
          console.log(data, 'setmanagerreq');
          if (data == null) {

          } else {
            this.assigningForm.patchValue
          }

        },
        error: (error: Error) => {
          console.error(error);

        }
      })

    });


    this.route.paramMap.subscribe(params => {
      this.employeeId = params.get('id');
      if (this.employeeId) {
        // Fetch employee data by ID
        this._employeedata.getbyid(this.employeeId).subscribe({
          next: (data: any) => {
            console.log(data, 'getdatatoassign');

            const employee = data.employee;

            if (employee) {
              this.assigningForm.patchValue({
                empId: employee.empId, // or however you retrieve empId from your data
                employee_id: employee.empId
              });
            }
          },
          error: (err: any) => {
            console.error('Error fetching employee data', err);
          }
        });
      }
    });


    this.assigningForm.get('role_id')?.valueChanges.subscribe(managerName =>
      this.dataService.getRolesById(managerName).subscribe({
        next: (data: any) => {
          console.log(data);
          this.roledata=data
          
          if (data.roleName !== "Manager") {
            this.assigningForm.get('tech_lead')?.enable()
          } else {
            this.assigningForm.get('tech_lead')?.disable()

          }
        }
      })
    )
    console.log(this.managerName,'this.managernAME');
    
  }

  
  getEmployeeAsManager(deptId: number): void {

      this.dataService.getEmployeeAsManager(deptId).subscribe({
        next: (data: any) => {
          console.log(data, 'data');
          if (data != null) {
            this._employeedata.getbyid(data.employee_id).subscribe({
              next: (data: any) => {
                this.managerName = `${data.firstName} ${data.lastName}`,
                  this.managerId = data.empId;
              
              }
            })
          } else if (data == null) {
            this.disableemployeeselect = true;
            this.managerName = '';
          }
        },
        error: (error: Error) => {
          console.error(error);

        }
      })
  }

  onSubmit() {
    console.log(this.assigningForm.value, 'values');

    // Create a new FormData object
    const formData = new FormData();

    // Append form values to the FormData object
    formData.append('employee_id', this.assigningForm.value.employee_id);
    formData.append('dept_id', this.assigningForm.value.dept_id);

console.log(this.roledata.role_id,'roledataid');

    if (this.assigningForm.value.role_id === this.roledata.role_id ) {
      formData.append('role_id', '0'); // Append '0' as a string
    } else {
      formData.append('role_id', this.assigningForm.value.role_id);
    }

    // Conditionally append the tech_lead value
    if (this.assigningForm.value.tech_lead === undefined) {
      formData.append('tech_lead', '0'); // Append '0' as a string
    } else {
      formData.append('tech_lead', this.assigningForm.value.tech_lead);
    }

    // Now you can use formData in your HTTP request or further processing

    console.log(formData, 'formData');

    const payload = {
      employee_id: this.assigningForm.value.employee_id,
      dept_id: this.assigningForm.value.dept_id,
      role_id: this.assigningForm.value.role_id,
      tech_lead: this.assigningForm.value.tech_lead === undefined ? 0 : this.assigningForm.value.tech_lead
  };

    this.dataService.postAsssignData(payload).subscribe(
      {
        next: (data: any) => {
          this.assignData = data;
          console.log(data, 'assigned succesfully');
          this.toaster.success("Submitted Successfully", "Success");
          this.router.navigate(['/hierarchichal-chart']);
        },
        error: (err: Error) => {
          console.error(err);
          this.toaster.error("User Already Assigned");

        },
        complete() {
          this._postdata.addmsg('New Employee Added: ');

        },
      }

    )
  }

}


