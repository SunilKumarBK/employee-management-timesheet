import { Routes } from '@angular/router';
import { EmployeeListComponent } from './employee-list/employee-list.component';
import { AddEmployeeComponent } from './add-employee/add-employee.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { OrganizationChartComponent } from './organization-chart/organization-chart.component';
import { AssignEmployeeComponent } from './assign-employee/assign-employee.component';
import { ProfileComponent } from './profile/profile.component';
import { DialogComponent } from './dialogBody/dialog/dialog.component';
import { NotfoundcomponentComponent } from './notfoundcomponent/notfoundcomponent.component';
import { Component } from '@angular/core';
import { TimesheetComponent } from './timesheet/timesheet.component';
import { RegisterComponent } from './register/register.component';
import { LoginComponent } from './login/login.component';
import { LogoutComponent } from './logout/logout.component';

export const routes: Routes = [
  { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
  {
    path: 'dashboard',
    // component:DashboardComponent,
    loadComponent: () =>
      import('./dashboard/dashboard.component').then(
        (m) => m.DashboardComponent
      ),
  },
  {
    path: 'employee',
    children: [
      {
        path: 'employeelist',
        component: EmployeeListComponent,
      },
      {
        path: 'profile',
        component: ProfileComponent,
      },
      {
        path: 'timesheet/:empId',
        component: TimesheetComponent,
      },
      {
        path: 'register',
        component: RegisterComponent,
      },
      {
        path: 'login',
        component: LoginComponent,
      },
      {
        path: 'logout',
        component: LogoutComponent,
      },
    ],
  },
  {
    path: 'addemployee',
    component: AddEmployeeComponent,
  },
  {
    path: 'addemployee/:id',
    component: AddEmployeeComponent,
  },
  {
    path: 'hierarchichal-chart',
    component: OrganizationChartComponent,
  },
  {
    path: 'assignemployee',
    component: AssignEmployeeComponent,
  },
  {
    path: 'assignemployee/:id',
    component: AssignEmployeeComponent,
  },
  {
    path: 'profile',
    component: ProfileComponent,
  },
  {
    path: '**',
    component: NotfoundcomponentComponent,
  },
];
