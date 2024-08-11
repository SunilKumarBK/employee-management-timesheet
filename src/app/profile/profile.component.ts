import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardContent, MatCardFooter, MatCardHeader, MatCardImage, MatCardModule } from '@angular/material/card';
import { EmployeeService } from '../service/employee.service';
import { HttpErrorResponse } from '@angular/common/http';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [MatCardModule, MatCardContent, MatCardHeader, MatCardFooter, MatCardImage, ReactiveFormsModule, CommonModule,MatIconModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {

  constructor(private _employeeservice: EmployeeService) { }

  profileForm: FormGroup;
  id: number;
  employee: any;
  data: any;
  aadharPdfUrl: string;
  getdocuments: any;

  ngOnInit(): void {
    this.profileForm = new FormGroup({
      firstName: new FormControl({ value: '', disabled: true }),
      lastName: new FormControl({ value: '', disabled: true }),
      email: new FormControl({ value: '', disabled: true }, [Validators.required, Validators.email]),
      phoneNo: new FormControl({ value: '', disabled: true }, [Validators.required]),
      fatherName: new FormControl({ value: '', disabled: true }),
      dateOfBirth: new FormControl({ value: '', disabled: true }),
      emergencyContact: new FormControl({ value: '', disabled: true }, [Validators.required]),
      address: new FormControl({ value: '', disabled: true }, [Validators.required]),
      gender: new FormControl('', Validators.required),
      relationstatus: new FormControl('', Validators.required),
      bloodgroup: new FormControl('', Validators.required),
      aadhar: new FormControl(null, Validators.required),
      profilephoto: new FormControl(null, Validators.required),
      empId: new FormControl('')
    });

    this.id = this._employeeservice.idtodialog;
    console.log(this.id, 'idd');

    // Fetch employee details
    this._employeeservice.getbyid(this.id).subscribe({
      next: (data: any) => {
        this.employee = data.employee;
        console.log(this.employee);
        this.profileForm.patchValue(data.employee);
      },
      error: (error: HttpErrorResponse) => {
        console.error('Error fetching employee details', error);
      }
    });

    // Fetch document details
    this._employeeservice.getDocuments(this.id).subscribe({
      next: (data: any) => {
        this.getdocuments = data;
        console.log(this.getdocuments, 'getdocuments');

        // Convert base64 to Blob and File
        const aadharBlob = this.convertBase64ToBlob(data.aadharImage, 'application/pdf');
        const aadharFile = this.convertBlobToFile(aadharBlob, 'aadhar.pdf');
        this.aadharPdfUrl = this.createObjectUrlFromBlob(aadharBlob);

        const profilePhotoBlob = this.convertBase64ToBlob(data.profilePhoto, 'image/jpeg');
        const profilePhotoFile = this.convertBlobToFile(profilePhotoBlob, 'profile.jpg');
        this.data = this.createObjectUrlFromBlob(profilePhotoBlob);

        // Patch form values
        this.profileForm.patchValue({
          gender: data.gender,
          relationstatus: data.relationship,
          bloodgroup: data.bloodgroup,
          aadhar: aadharFile,
          profilephoto: profilePhotoFile
        });
        console.log(data.gender, data.relationship, data.bloodgroup, data.aadharImage, this.data);
      },
      error: (err) => {
        console.error('Error fetching documents', err);
      }
    });
  }

  // Convert Base64 string to Blob
  convertBase64ToBlob(base64: string, contentType: string): Blob {
    const byteCharacters = atob(base64);
    const byteNumbers = new Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
      byteNumbers[i] = byteCharacters.charCodeAt(i);
    }
    const byteArray = new Uint8Array(byteNumbers);
    return new Blob([byteArray], { type: contentType });
  }

  // Convert Blob to File
  convertBlobToFile(blob: Blob, fileName: string): File {
    return new File([blob], fileName, { type: blob.type });
  }

  // Create Object URL from Blob
  createObjectUrlFromBlob(blob: Blob): string {
    return URL.createObjectURL(blob);
  }

  // Handle file input change
  onFileChange(event: any, type: string): void {
    const file = event.target.files[0];
    if (file) {
      if (type === 'aadhar') {
        this.profileForm.patchValue({ aadhar: file });
      } else if (type === 'profilephoto') {
        this.profileForm.patchValue({ profilephoto: file });
      }
    }
  }

  onImagePicked(event: Event, controlName: string) {
    const file = (event.target as HTMLInputElement).files[0]; // Here we use only the first file (single file)
    // this.profileForm.patchValue({ document:{controlName: file}});
    if (file) {
      this.profileForm.patchValue({
        // document:{
        [controlName]: file
        // }
      });
    }
  }



  // Download Aadhar PDF
  downloadAadharPdf() {
    const link = document.createElement('a');
    link.href = this.aadharPdfUrl;
    link.download = 'aadhar.pdf';
    link.click();
  }

  // Update employee details
  onUpdate() {
    console.log('toupdate');
    alert(111);

    const formData = new FormData();
    formData.append('empId', this.profileForm.value.empId);

    // Append JSON data as a string
    formData.append('personalDetails', JSON.stringify({
      empId: this.profileForm.value.empId,
      gender: this.profileForm.value.gender || this.getdocuments.gender,
      relationship: this.profileForm.value.relationstatus || this.getdocuments.relationship,
      bloodgroup: this.profileForm.value.bloodgroup || this.getdocuments.bloodgroup,
    }));

    const aadharFile = this.profileForm.get('aadhar').value;
    const profilePhoto = this.profileForm.get('profilephoto').value;

    // Append files to FormData
    if (aadharFile instanceof File) {
      formData.append('aadhar', aadharFile);
    } else {
      // Reuse existing file if not a new one
      const aadharBlob = this.convertBase64ToBlob(this.getdocuments.aadharImage, 'application/pdf');
      formData.append('aadhar', this.convertBlobToFile(aadharBlob, 'aadhar.pdf'));
    }

    if (profilePhoto instanceof File) {
      formData.append('profilephoto', profilePhoto);
    } else {
      // Reuse existing file if not a new one
      const profilePhotoBlob = this.convertBase64ToBlob(this.getdocuments.profilePhoto, 'image/jpeg');
      formData.append('profilephoto', this.convertBlobToFile(profilePhotoBlob, 'profile.jpg'));
    }

    alert('before');
    this._employeeservice.updateDocs(formData).subscribe({
      next: (data: any) => {
        console.log(data, 'successfully updated');
      },
      error: (error: HttpErrorResponse) => {
        console.log('update failed');
        console.error(error);
      }
    });
    alert('after');
  }

  // Submit form data
  onSubmit() {
    alert('submit');
    console.log(this.profileForm.value);

    const formData = new FormData();
    formData.append('empId', this.profileForm.value.empId);

    // Append JSON data as a string
    formData.append('personalDetails', JSON.stringify({
      empId: this.profileForm.value.empId,
      gender: this.profileForm.value.gender,
      relationship: this.profileForm.value.relationstatus,
      bloodgroup: this.profileForm.value.bloodgroup,
      aadhar: this.profileForm.value.aadhar,
      profilephoto: this.profileForm.value.profilephoto
    }));

    const aadharFile = this.profileForm.get('aadhar').value;
    const profilePhoto = this.profileForm.get('profilephoto').value;

    if (aadharFile instanceof File) {
      formData.append('aadhar', aadharFile);
    }

    if (profilePhoto instanceof File) {
      formData.append('profilephoto', profilePhoto);
    }
    console.log(formData, 'formdata');

    if (this.profileForm.value.aadhar && this.profileForm.value.profilephoto) {
      console.log(this.profileForm.value);

      this._employeeservice.uploadDocuments(formData)
        .subscribe({
          next: (response: any) => {
            console.log('Files uploaded successfully', response);
          },
          error: (error: HttpErrorResponse) => {
            console.error('Upload failed', error);
          }
        });
    }
  }
}
