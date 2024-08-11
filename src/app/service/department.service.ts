import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class DepartmentService {


  private apiUrl = 'http://localhost:8080'; // Your Go API endpoint

  constructor(private http: HttpClient) { }

  getRolesByDepartment(deptId: number): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/getrolesbydepartment/${deptId}`);
  }

  getDepartments(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/departments`);
  }

  getDepartmentsById(deptId: number): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/departmentbyid/${deptId}`);
  }

  getManagers(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/managers`);
  }
  getManagerByDepartment(deptId: number): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/manager/${deptId}`);
  }
   //post
   postAsssignData(value: any): Observable<[]> {
    console.log(value);
    
    return this.http.post<[]>(`${this.apiUrl}/assignemployee`, value)
  }

  //get joined data of department employee role manager for chart
  getHierarchyData():Observable<any[]>{
    return this.http.get<any>(`${this.apiUrl}/hierarchy`)
  }
  getHierarchyDatabyid(id:number):Observable<any[]>{
    return this.http.get<any>(`${this.apiUrl}/hierarchy/${id}`)
  }
  getEmployeeAsManager(id:number):Observable<any[]>{
    return this.http.get<any>(`${this.apiUrl}/getEmployeeAsManager/${id}`)
  }

  getRolesById(deptId: number): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/getRoleById/${deptId}`);
  }

  getAssign():Observable<any[]>{
    return this.http.get<any[]>(`${this.apiUrl}/getAssignData`);
  }
}
