import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators ,ReactiveFormsModule, FormControl} from '@angular/forms';
import { AuthService } from '../service/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css'],
  standalone:true,
  imports:[ReactiveFormsModule]
})
export class RegisterComponent {
  registerForm: FormGroup;
  errorMessage: string = '';

  constructor(private fb: FormBuilder, private authService: AuthService, private router: Router) {
    this.registerForm = new FormGroup({
      firstname:new FormControl ('', Validators.required),
      lastname:new FormControl ('', Validators.required),
      email:new FormControl ('', [Validators.required, Validators.email]),
      password:new FormControl ('', Validators.required)
    });
  }

  onSubmit() {
    alert(0)
    // if (this.registerForm.valid) {
      this.authService.register(this.registerForm.value).subscribe({
        next:(response) => {
          this.router.navigate(['/login']);
        },
        error:(error) => {
          this.errorMessage = 'Registration failed';
        }
  });
    // }
  }
}
