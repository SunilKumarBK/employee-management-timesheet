import { Component, inject, output, ViewChild } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon'
import { CommonModule } from '@angular/common';
import { EmployeeService } from '../service/employee.service';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input'
import { MatFormFieldModule } from '@angular/material/form-field'
import { FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { FilterPipe } from "../filter.pipe";
import { NgxSpinnerModule } from 'ngx-spinner';
import { ToastrService } from 'ngx-toastr';
import { DialogComponent } from '../dialogBody/dialog/dialog.component';
import { MatDialog, MatDialogActions, MatDialogClose, MatDialogModule, MatDialogContent } from '@angular/material/dialog';
import { MatPaginatorModule } from '@angular/material/paginator';
import { IAddemployeeData } from '../addemployeedata';
import { PaginationComponent } from "../pagination/pagination.component";
import {MatRippleModule} from '@angular/material/core';
import { DashboardComponent } from '../dashboard/dashboard.component';



@Component({
  selector: 'app-employee-list',
  standalone: true,
  templateUrl: './employee-list.component.html',
  styleUrl: './employee-list.component.css',
  imports: [
    RouterLink,
    RouterLinkActive,
    MatToolbarModule,
    MatIconModule,
    CommonModule,
    MatButtonModule,
    MatInputModule,
    MatFormFieldModule,
    FormsModule,
    ReactiveFormsModule,
    FilterPipe,
    NgxSpinnerModule,
    DialogComponent,
    MatDialogModule,
    MatDialogActions,
    MatDialogClose,
    MatDialogContent,
    MatPaginatorModule,
    PaginationComponent,
    MatRippleModule
]
})
export class EmployeeListComponent {

  toaster = inject(ToastrService);
  $event: any;
  constructor(private _employeedata: EmployeeService, private matDialog: MatDialog,private notification:DashboardComponent) { }

  title = ['S.no', 'Employee ID', 'Employee Name', 'Mail ID', 'Phone No', 'Address', 'Action'];
  itemsPerPage = 5;
  currentPage = 1;

  employeedata: IAddemployeeData[] = [];
  errorMsg: any;
  filtervalue: any = '';
  employeedetails: any;
  isDialog = false;
  deletDialog = false;
  deleteId: any;
  getId: any;
  getbyiddata: any[] = [];
  show = false;
  filterData: any[] = [];



  toggleDialog() {
    this.show = true;
  }

  ngOnInit() {

    this.reloadData();
  }

  passid: EmployeeService = inject(EmployeeService)

  openDialog(id: any) {
    console.log(id);

    this.passid.onPass(id);
    this.matDialog.open(DialogComponent, {
      width: '400px',
    })
  }




  reloadData() {
    this._employeedata.getUsersList().subscribe(
      {
        next: (data: any[]) => {
          this.employeedata = data;
        },
        error: (error) => {
          console.error('There was an error!', error);
        },
        complete: () => {
          console.log('Request complete');
        }
      });
   

  }

  delete(id: any) {
    console.log(id, 'deleteid');

    this._employeedata.delete(id).subscribe(
      {
        next: (data: any[]) => {
          // this.employeedata = data;
          this.toaster.success("Deleted Successfully", "Success");
          this.deletDialog = false;
          this.reloadData();
        },
        error: (error) => {
          console.error('There was an error!', error);
          this.toaster.error("Something went wrong", "Error")
        },
        complete: () => {
          console.log('Request complete');
          this._employeedata.delmssg('Employee deleted Sucessfully : ')

        }
      }

    );
  }

  getbyid(value: any) {
    this.isDialog = true;
    this.getId = value.empId;
    this._employeedata.delmsg(value);
    console.log(value, 'valueid');

    this._employeedata.getbyid(this.getId).subscribe({
      next: (data: any) => {
        this.getbyiddata = data;
        console.log(this.getbyiddata, 'getbydata'); // Log received data for debugging
      },
      error: (error: any) => {
        console.error('There was an error!', error);
      },
      complete: () => {
        console.log('Request complete');
      }
    });
    console.log(this.getbyiddata, 'getbyiddata');




  }


  deleteDialog(value: any) {
    console.log(value);
    this.deletDialog = true;
    this.deleteId = value;

  }

  closeDeleteDialog() {
    this.deletDialog = false;

  }

  get paginatedData() {
    const start = (this.currentPage - 1) * (this.itemsPerPage);
    const end = start + this.itemsPerPage;

    return this.employeedata.slice(start, end);
  }

  changePage(page: any) {
    this.currentPage = page;

  }







}




