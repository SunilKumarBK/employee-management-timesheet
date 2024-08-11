import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AssignEmployeeComponent } from './assign-employee.component';

describe('AssignEmployeeComponent', () => {
  let component: AssignEmployeeComponent;
  let fixture: ComponentFixture<AssignEmployeeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AssignEmployeeComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AssignEmployeeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
