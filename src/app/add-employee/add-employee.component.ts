import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormsModule, ReactiveFormsModule, FormGroup, FormControl, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { CommonModule } from '@angular/common';
import moment from 'moment';
import { EmployeeService } from '../service/employee.service';
import { ActivatedRoute, Router } from '@angular/router';
import { NgxSpinnerModule } from 'ngx-spinner';
import { ToastrService } from 'ngx-toastr';
import { MatRadioModule } from '@angular/material/radio';
import { DashboardComponent } from '../dashboard/dashboard.component';


@Component({
  selector: 'app-add-employee',
  standalone: true,
  imports: [FormsModule, MatCardModule, MatInputModule, MatButtonModule, CommonModule, ReactiveFormsModule, NgxSpinnerModule, MatRadioModule],
  templateUrl: './add-employee.component.html',
  styleUrl: './add-employee.component.css'
})
export class AddEmployeeComponent implements OnInit {


  handleError: any;
  constructor(private router: Router, private _postdata: EmployeeService, private fb: FormBuilder, private route: ActivatedRoute, private notification: DashboardComponent) { }

  postdata: any;
  companyData: any;
  postdataerr: any;
  displayy = false;
  showxtra = false;
  employeeId: any;
  isEditMode: boolean = false;
  toaster = inject(ToastrService);
  radioOptions: string[] = ['Yes', 'No']
  getbyiddata: any;


  employeeForm: FormGroup;
  secondcmpForm: FormGroup;
  ngOnInit() {

    this.employeeForm = new FormGroup({
      firstName: new FormControl('', [Validators.required, Validators.minLength(4)]),
      lastName: new FormControl(''),
      email: new FormControl('', [Validators.required, Validators.email]),
      phoneNo: new FormControl('', [Validators.required]),
      fatherName: new FormControl(''),
      dateOfBirth: new FormControl(''),
      emergencyContact: new FormControl('', [Validators.required]),
      address: new FormControl('', [Validators.required]),
      empId: new FormControl({ value: '', disabled: false }, [Validators.required]),
      qualification: new FormControl('', [Validators.required]),
      experience: new FormControl({ value: '', disabled: false }, [Validators.required]),
      companyName: new FormControl('', [Validators.required, Validators.minLength(4)]),
      designation: new FormControl('', [Validators.required]),
      joinDate: new FormControl('', [Validators.required]),
      relievedDate: new FormControl('', [Validators.required]),
      totalDuration: new FormControl(''),
    });

    this.secondcmpForm = new FormGroup({
      companyName: new FormControl('', [Validators.required, Validators.minLength(4)]),
      designation: new FormControl('', [Validators.required]),
      joinDate: new FormControl('', [Validators.required]),
      relievedDate: new FormControl('', [Validators.required]),
      duration: new FormControl('')

    })




    // Subscribe to joinDate and relievedDate changes to update totalDuration
    this.employeeForm.get('joinDate').valueChanges.subscribe(() => {
      this.updateTotalDuration();
    });

    this.employeeForm.get('relievedDate').valueChanges.subscribe(() => {
      this.updateTotalDuration();
    });

    this.secondcmpForm.get('joinDate').valueChanges.subscribe(() => {
      this.updateTotalDuration();
    });

    this.secondcmpForm.get('relievedDate').valueChanges.subscribe(() => {
      this.updateTotalDuration();
    });


    this.route.paramMap.subscribe(params => {
      this.employeeId = params.get('id');
      if (this.employeeId) {
        this.isEditMode = true;
        this.getEmployeeData(this.employeeId);
      }
    });



    // Patch the form data
    this.employeeForm.patchValue({ experience: 'false' });
  }


  isSubmitDisabled(): boolean {
    const experience = this.employeeForm.get('experience').value === "true";
    const initialFieldsInvalid = this.employeeForm.get('firstName').invalid ||
      this.employeeForm.get('email').invalid ||
      this.employeeForm.get('phoneNo').invalid ||
      this.employeeForm.get('emergencyContact').invalid ||
      this.employeeForm.get('address').invalid ||
      this.employeeForm.get('empId').invalid ||
      this.employeeForm.get('qualification').invalid;

    if (!experience) {
      return initialFieldsInvalid;
    } else {
      const companyFieldsInvalid = this.employeeForm.get('firstName').invalid ||
        this.employeeForm.get('email').invalid ||
        this.employeeForm.get('phoneNo').invalid ||
        this.employeeForm.get('emergencyContact').invalid ||
        this.employeeForm.get('address').invalid ||
        this.employeeForm.get('empId').invalid ||
        this.employeeForm.get('companyName').invalid ||
        this.employeeForm.get('designation').invalid ||
        this.employeeForm.get('joinDate').invalid ||
        this.employeeForm.get('relievedDate').invalid;
      return companyFieldsInvalid;
    }
  }






  //url var
  _url: string = '/assets/data/postdata.json';



  //post
  onClick() {
    // const formValue = this.employeeForm.value+this.secondcmpForm;
    const employeeFormValue = this.employeeForm.value;

    // If you need to include the second company form values
    const secondCompanyFormValue = this.secondcmpForm.value;

    const formValue = this.employeeForm.value;
    console.log(formValue, 'value');
    console.log(this.secondcmpForm.value, 'secondformvalue');


    formValue.experience = formValue.experience === "true";
    console.log(this.employeeForm.value, 'postvalues')
    const formData = {
      ...employeeFormValue,
      secondCompanyFormValue: secondCompanyFormValue
    };
    console.log(formData, 'formData');

    this._postdata.postUserdata(formData).subscribe(
      {
        next: (data: any) => {
          this.postdata = data;
          console.log(this.postdata, 'empolyeedata');
          this.toaster.success("Submitted Successfully", "Success");
          this.router.navigate(['/profile']);

        },
        error: (err: Error) => {
          this.postdataerr = err.message;
          console.log(err);
          this.toaster.error("Something went wrong", "Error");
        },
        complete: () => {
          console.log('New Employee Added.');
          this._postdata.addmsg('New Employee Added: ');

        }
      }
    );





  }



  //update

  onUpdate(): void {
    const experience = this.employeeForm.get('experience')?.value == 'true';
    const formData = {
      ...this.employeeForm.value, secondCompanyFormValue: this.secondcmpForm.value,
      experience: experience
    };
    console.log(formData, 'updatedata');


    this._postdata.updateEmployee(this.employeeId, formData).subscribe(
      {
        next: (data: any) => {
          console.log('Employee updated successfully:', data);
          this.toaster.success("Updated Successfully", "Success");
          this.router.navigate(['/dashboard']); // Navigate to list after updating
        },
        error: (error: any) => {
          console.error('Error updating employee:', error);
          this.toaster.error("Something went wrong", "Error");
        }
      });
  }


  //getbyid

  getEmployeeData(id: any): void {
    this._postdata.getbyid(id).subscribe(
      {
        next: (data: any) => {
          this.getbyiddata = data;
          this.employeeForm.patchValue(data.employee);
          const experienceValue = data.employee.experience ? 'true' : 'false';

          this.employeeForm.patchValue({ experience: experienceValue });

          // this.employeeForm.patchValue({ experience: data.employee.experience });



          console.log(this.employeeForm.value);
          console.log(data, 'getbyiddata');


          if ('experience' in data.employee) {
            this.employeeForm.patchValue({ experience: experienceValue });
            if (data.employee.experience === true) {
              this.displayy = true;
              this.employeeForm.patchValue(data.company[0]);
              if (data.company[1]) {
                this.showxtra = true;
                this.secondcmpForm.patchValue(data.company[1]);

              }
            }
          } else {
            this.employeeForm.patchValue({ experience: experienceValue }); // Default to false if not present
          }




          // If there's company data, patch it as well
          // if (data.company) {
          //   this.employeeForm.patchValue(data.company);
          // }
        },
        error: (error: any) => {
          console.error('Error fetching employee data:', error);
        }
      });
    this.employeeForm.get('experience')?.valueChanges.subscribe((data: any) => {
      const exp = data;
      if (data == 'true') {
        this.displayy = true;
        this.employeeForm.patchValue(this.getbyiddata.company[0]);
        if (this.getbyiddata.company[1]) {
          this.showxtra = true;
          this.secondcmpForm.patchValue(this.getbyiddata.company[1]);

        }
      }
    })




    // this.employeeForm.get('experience')?.disable()
    this.employeeForm.get('empId')?.disable()
  }


  experience(): void {
    this.displayy = true;
    if (this.employeeForm.get('experience').value === "true") {
      this.employeeForm.patchValue({ experience: "false" });
    } else {
      this.employeeForm.patchValue({ experience: "true" });
    }
  }

  experiencee(): void {
    this.displayy = false;
    if (this.employeeForm.get('experience').value === 'false') {
      this.employeeForm.patchValue({ experience: 'true' });
    } else {
      this.employeeForm.patchValue({ experience: 'false' });
    }
  }

  togglebtn() {
    this.showxtra = true
  }
  togglebtnremove() {
    this.showxtra = false;
  }

  // startDate: string | any;
  // endDate: string | any;

  calculateDuration(startDate: string, endDate: string): string {
    const start = moment(startDate);
    const end = moment(endDate);
    const duration = moment.duration(end.diff(start));
    const years = duration.years();
    const months = duration.months();
    const days = duration.days();

    let durationString = '';
    if (years > 0) {
      durationString += years + ' years ';
    }
    if (months > 0) {
      durationString += months + ' months ';
    }
    if (days > 0) {
      durationString += days + ' days';
    }

    return durationString.trim();

  }
  updateTotalDuration(): void {
    const startDate = this.employeeForm.get('joinDate').value;
    const endDate = this.employeeForm.get('relievedDate').value;
    const duration = this.calculateDuration(startDate, endDate);



    // Patch the calculated duration into totalDuration control
    this.employeeForm.patchValue({
      totalDuration: duration
    });


    const seccmpstartDate = this.secondcmpForm.get('joinDate').value;
    const seccmpendDate = this.secondcmpForm.get('relievedDate').value;
    const seccduration = this.calculateDuration(seccmpstartDate, seccmpendDate);

    this.secondcmpForm.patchValue({
      duration: seccduration
    });

  }


}


