import { Component, OnInit } from '@angular/core';
import { EmployeeService } from '../service/employee.service';
import { CommonModule } from '@angular/common';
import { Router, RouterLink, RouterLinkActive } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive,MatIconModule],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.css'
})
export class SidebarComponent implements OnInit {
  collapse: any;
  isCollapsed: boolean = false;
  constructor(private isCollapse: EmployeeService) { }

  ngOnInit() {
    this.isCollapse.currentData.subscribe(data => this.collapse = data);
   
  }

  sidebar() {
    this.isCollapsed = !this.isCollapsed;
    this.isCollapse.changeData(this.isCollapsed);
  }



}
