import { CommonModule } from '@angular/common';
import { Component, Input, OnInit, OnChanges, SimpleChanges, Output } from '@angular/core';
import { EventEmitter } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-pagination',
  standalone: true,
  imports: [CommonModule ,MatIconModule],
  templateUrl: './pagination.component.html',
  styleUrl: './pagination.component.css'
})
export class PaginationComponent implements OnInit, OnChanges {
  @Input() totalItems: any;
  @Input() currentPage: number;
  @Input() itemsPerPage: number;
  @Output() onClick = new EventEmitter<any>();

  totalPages = 0;
  pages:number[] =[]

  constructor() {}

  ngOnInit(): void {
    // Initial calculation if the inputs are already available
    this.calculateTotalPages();
    console.log(this.totalItems,'totalitems');

  }

  ngOnChanges(changes: SimpleChanges): void {
    // Handle changes to the inputs
    if (changes['totalItems'] || changes['itemsPerPage']) {
      this.calculateTotalPages();
    }
  }

  private calculateTotalPages(): void {
    if (this.totalItems && this.itemsPerPage) {
      this.totalPages = Math.ceil(this.totalItems / this.itemsPerPage);
      console.log(`Total pages: ${this.totalPages}`);
      this.pages = Array.from({length: this.totalPages},(_,i)=>i + 1);
    }
  }


  pageClicked(page:any):void{
    if(page>this.totalPages) return;
    this.onClick.emit(page);
    

  }


}
