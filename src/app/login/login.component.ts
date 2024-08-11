import { Component } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { AuthService } from '../service/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  standalone:true,
  imports:[ReactiveFormsModule]
})
export class LoginComponent {
  loginForm: FormGroup;
  errorMessage: string = '';

  constructor(private fb: FormBuilder, private authService: AuthService, private router: Router) {
    this.loginForm = new FormGroup({
      email:new FormControl ('', [Validators.required, Validators.email]),
      password:new FormControl ('', Validators.required)
    });
  }

  onSubmit() {
    alert(1)
    // if (this.loginForm.valid) {
      this.authService.login(this.loginForm.value).subscribe({
        next:(response:any) => {
          console.log(response,'login succesfully');
          
          // this.router.navigate(['employee/timesheet']);
        },
        error:(error:Error) => {
          this.errorMessage = 'Invalid email or password';
        }
  });
    // }
  }
}
