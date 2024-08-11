import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private baseUrl = 'http://localhost:8080'; // Adjust according to your API base URL

  constructor(private http: HttpClient) {}

  register(credentials: any): Observable<any> {
    return this.http.post(`${this.baseUrl}/register`, credentials);
  }

  login(credentials: any): Observable<any> {
    return this.http.post(`${this.baseUrl}/login`, credentials);
  }

  logout(id:any): Observable<any> {
    return this.http.post<any[]>(`${this.baseUrl}/logout`,{empId:id});
  }

  getTimesheet(empId:number): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/timesheet/${empId}`);
  }
}
