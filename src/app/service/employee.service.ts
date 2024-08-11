import { EventEmitter, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { IEmployeeData } from '../employeeData';
import { HttpClient } from '@angular/common/http';
import { IAddemployeeData } from '../addemployeedata';
import { BehaviorSubject } from 'rxjs';
import { AddEmployeeComponent } from '../add-employee/add-employee.component';
import { DashboardComponent } from '../dashboard/dashboard.component';

@Injectable({
  providedIn: 'root'
})
export class EmployeeService {
  private baseUrl = 'http://localhost:8080';

  constructor(private _http: HttpClient) { }

  private dataSource = new BehaviorSubject<any>(null);
  currentData = this.dataSource.asObservable();

  addedmessage:any='';
  deletedmessage:any;
  dellmsg;


  addmsg(msg:any){
this.addedmessage=msg;
console.log(this.addedmessage,'thisaddedmsg');


  }

  delmsg(msg:any){
    this.deletedmessage=msg.firstName+msg.lastName;
    console.log(this.deletedmessage,'deleted message');
    
  }

  delmssg(msg:any){
    this.dellmsg=msg;
  }



  changeData(data: any) {
    this.dataSource.next(data);

  }

  idtodialog: any;
  passid: EventEmitter<any> = new EventEmitter<any>();

  onPass(value: any) {
    console.log(value);

    this.idtodialog = value;
    this.passid.emit(value);

  }




  //get
  getUsersList(): Observable<any[]> {

    return this._http.get<any[]>(`${this.baseUrl}/employee`);
  }

  //getbyid
  getbyid(id: number): Observable<any[]> {
    return this._http.get<any[]>(`${this.baseUrl}/getbyid/emplywithprevcompany/${id}`);
  }

  //post
  postUserdata(value: any): Observable<IAddemployeeData[]> {
    return this._http.post<IAddemployeeData[]>(`${this.baseUrl}/addemployee`, value)
  }

  //delete
  delete(id: any): Observable<any[]> {
    return this._http.delete<any[]>(`${this.baseUrl}/delete/employee/${id}`)
  }


  //update
  updateEmployee(id: string, data: any): Observable<any> {
    return this._http.put(`${this.baseUrl}/update/employee/${id}`, data);
  }


  uploadDocuments(data:any): Observable<any> {
    console.log(data,'adadd');
    
    const formData = new FormData();
    formData.append('aadhar', data.aadhar);
    formData.append('profilephoto', data.profilephoto);
    formData.append('empId', data.empId);
    console.log(formData,'formdatad');
    

    return this._http.post(`${this.baseUrl}/handlePersonalDetailsAndDocuments`, data);
  }


  getDocuments(id: number): Observable<any[]> {
    return this._http.get<any[]>(`${this.baseUrl}/getDocuments/${id}`);
  }

  updateDocs( data: any): Observable<any> {
    return this._http.put(`${this.baseUrl}/handleUpdatePersonalDetailsAndDocuments`, data);
  }

  setManager(id:number):Observable<any>{
    return this._http.get(`${this.baseUrl}/setManager/${id}`);
  }
}
