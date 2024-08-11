import { Component, OnInit } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { NavbarComponent } from './navbar/navbar.component';
import { NgxSpinnerModule } from 'ngx-spinner'
import { SidebarComponent } from './sidebar/sidebar.component';
import { EmployeeService } from './service/employee.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet,NavbarComponent,NgxSpinnerModule,SidebarComponent,CommonModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  title = 'employee-management';
  collapse:any;
  displaynone:boolean=false;


  constructor(private isCollapse: EmployeeService,private router:Router) { }


  ngOnInit():void{
    this.isCollapse.currentData.subscribe(data => this.collapse = data);
    this.router.events.subscribe(() => {
      this.checkroute();
    });
    this.checkroute();
  }
  checkroute(){
    const checkRoute = this.router.url;
    this.displaynone=checkRoute.includes('hierarchichal-chart');
    
      }
}
