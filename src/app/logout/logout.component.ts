import { Component } from '@angular/core';
import { AuthService } from '../service/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-logout',
  templateUrl: './logout.component.html',
  styleUrls: ['./logout.component.css']
})
export class LogoutComponent {
  id:any=1;

  constructor(private authService: AuthService, private router: Router) {}

  onLogout() {
    const formData= new FormData()

    formData.append('empId','1');
    this.authService.logout(this.id).subscribe({
      next:(response:any) => {
        console.log('Logout successful', response);
        this.router.navigate(['/login']);  // Navigate to the login page after logout
      },
      error:(error:Error) => {
        console.error('Logout failed', error);
      }
   } );
  }
}
