package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"fmt"
    "log"
    "net/http"
    "time"

	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	  "github.com/gorilla/sessions"
    "golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

type CompanyDetails struct {
	CompanyName  string `json:"companyName"`
	Designation  string `json:"designation"`
	JoinDate     string `json:"joinDate"`
	RelievedDate string `json:"relievedDate"`
	Duration     string `json:"duration"`
	id           int    `json:"id`
}

type Employee struct {
	ID                     int            `json:"id"`
	EmpId                  int            `json:"empId"`
	FirstName              string         `json:"firstName"`
	LastName               string         `json:"lastName"`
	Email                  string         `json:"email"`
	PhoneNo                int            `json:"phoneNo"`
	FatherName             string         `json:"fatherName"`
	EmergencyContact       int            `json:"emergencyContact"`
	DateOfBirth            string         `json:"dateOfBirth"`
	Address                string         `json:"address"`
	Qualification          string         `json:"qualification"`
	Experience             bool           `json:"experience"`
	CreatedTime            string         `json:"created_time"`
	CompanyName            string         `json:"companyName"`
	Designation            string         `json:"designation"`
	JoinDate               string         `json:"joinDate"`
	RelievedDate           string         `json:"relievedDate"`
	TotalDuration          string         `json:"totalDuration"`
	SecondCompanyFormValue CompanyDetails `json:"secondCompanyFormValue"`
}

type EmployeeWithCompany struct {
	Employee Employee         `json:"employee"`
	Company  []CompanyDetails `json:"company"`
}

//

//get

func dataHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(` SELECT empId,firstName,lastName,email,phoneNo,fatherName,emergencyContact,dateOfBirth,address,experience,qualification,created_time FROM emply`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.EmpId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.PhoneNo, &emp.FatherName, &emp.EmergencyContact, &emp.DateOfBirth, &emp.Address, &emp.Experience, &emp.Qualification, &emp.CreatedTime); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved data: %v", employees)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//
//

// func prevcompanybyid(w http.ResponseWriter, r *http.Request) {
// 	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
// 	if err != nil {
// 		log.Printf("Error opening database: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	params := mux.Vars(r)
// 	idStr := params["id"]

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		log.Printf("Invalid employee ID: %v", err)
// 		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
// 		return
// 	}
// 	row, err := db.Query(` SELECT id,companyName,position,startDate,endDate,duration FROM prevcompany where empId=?`, id)
// 	if err != nil {
// 		log.Printf("Error executing query: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	defer row.Close()

// 	var company []CompanyDetails
// 	for row.Next() {
// 		var emp CompanyDetails
// 		if err := row.Scan(&emp.id, &emp.CompanyName, &emp.Designation, &emp.JoinDate, &emp.RelievedDate, &emp.Duration); err != nil {
// 			log.Printf("Error scanning row: %v", err)
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		company = append(company, emp)
// 	}

// 	if err := row.Err(); err != nil {
// 		log.Printf("Rows error: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	log.Printf("Retrieved company data: %v", company)

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(company); err != nil {
// 		log.Printf("Error encoding response: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func prevcompanybyid(tx *sql.Tx, empId string) ([]CompanyDetails, error) {
	rows, err := tx.Query(`SELECT id, companyName, position, startDate, endDate, duration FROM prevcompany WHERE empId=?`, empId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []CompanyDetails
	for rows.Next() {
		var company CompanyDetails
		if err := rows.Scan(&company.id, &company.CompanyName, &company.Designation, &company.JoinDate, &company.RelievedDate, &company.Duration); err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}

func emplybyid(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}
	row, err := db.Query(` SELECT id,empId,firstName,lastName,email,phoneNo,fatherName,emergencyContact,dateOfBirth,address,experience,qualification FROM emply where empId=?`, id)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer row.Close()

	var emply []Employee
	for row.Next() {
		var emp Employee
		if err := row.Scan(&emp.ID, &emp.EmpId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.PhoneNo, &emp.FatherName, &emp.EmergencyContact, &emp.DateOfBirth, &emp.Address, &emp.Experience, &emp.Qualification); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		emply = append(emply, emp)
	}

	if err := row.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved company data: %v", emply)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emply); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//
//

//employee and company by id

func employeeWithCompanyById(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var result EmployeeWithCompany

	// Query for employee details
	err = db.QueryRow("SELECT empId,firstName,lastName,email,phoneNo,fatherName,emergencyContact,dateOfBirth,address,experience,qualification,created_time FROM emply WHERE empId = ?", id).Scan(
		&result.Employee.EmpId,
		&result.Employee.FirstName,
		&result.Employee.LastName,
		&result.Employee.Email,
		&result.Employee.PhoneNo,
		&result.Employee.FatherName,
		&result.Employee.EmergencyContact,
		&result.Employee.DateOfBirth,
		&result.Employee.Address,
		&result.Employee.Experience,
		&result.Employee.Qualification,
		&result.Employee.CreatedTime,
	)
	if err != nil {
		log.Printf("Error executing query for employee: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query for previous company details
	rows, err := db.Query("SELECT companyName, position, startDate, endDate, duration FROM prevcompany WHERE empId = ?", id)
	if err != nil {
		log.Printf("Error executing query for previous company: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var company CompanyDetails
		if err := rows.Scan(&company.CompanyName, &company.Designation, &company.JoinDate, &company.RelievedDate, &company.Duration); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result.Company = append(result.Company, company)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved employee with company data: %v", result)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//

//add

func addEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Optionally, validate or process emp data here

	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback if any error occurs before commit

	stmt, err := tx.Prepare("INSERT INTO emply (empId, firstName, lastName, email, phoneNo, fatherName, emergencyContact, dateOfBirth, address, qualification,experience) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute insert into emply table
	_, err = stmt.Exec(emp.EmpId, emp.FirstName, emp.LastName, emp.Email, emp.PhoneNo, emp.FatherName, emp.EmergencyContact, emp.DateOfBirth, emp.Address, emp.Qualification, emp.Experience)
	if err != nil {
		log.Printf("Error executing insert statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if emp.Experience {
		// Insert into prevcompany table
		companyStmt, err := tx.Prepare("INSERT INTO prevcompany (companyName, position, startDate, endDate, duration, empId) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("Error preparing company statement: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer companyStmt.Close()

		_, err = companyStmt.Exec(emp.CompanyName, emp.Designation, emp.JoinDate, emp.RelievedDate, emp.TotalDuration, emp.EmpId)
		if err != nil {
			log.Printf("Error executing company insert statement: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// log.Printf(_)
		// Insert into second company table if needed
		if emp.SecondCompanyFormValue.CompanyName != "" {
			_, err = companyStmt.Exec(emp.SecondCompanyFormValue.CompanyName, emp.SecondCompanyFormValue.Designation, emp.SecondCompanyFormValue.JoinDate, emp.SecondCompanyFormValue.RelievedDate, emp.SecondCompanyFormValue.Duration, emp.EmpId)
			if err != nil {
				log.Printf("Error executing second company insert statement: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Employee added successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

//delete

func deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Delete from prevcompany table first
	_, err = db.Exec("DELETE FROM prevcompany WHERE empId = ?", id)
	if err != nil {
		log.Printf("Error deleting from prevcompany: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM employeeRole WHERE employee_id = ?", id)
	if err != nil {
		log.Printf("Error deleting from prevcompany: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("DELETE FROM emply WHERE empId = ?")
	if err != nil {
		log.Printf("Error preparing delete statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		log.Printf("Error executing delete statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Printf("No employee found with ID: %d", id)
		http.Error(w, "No employee found with the given ID", http.StatusNotFound)
		return
	}

	log.Printf("Employee deleted successfully: ID %d", id)
	w.WriteHeader(http.StatusNoContent)
}

//update

func updateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var emp Employee
	// var company CompanyDetails

	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE emply SET firstName=?, lastName=?, email=?, phoneNo=?, fatherName=?, emergencyContact=?, dateOfBirth=?, address=?, experience=? WHERE empId=?", emp.FirstName, emp.LastName, emp.Email, emp.PhoneNo, emp.FatherName, emp.EmergencyContact, emp.DateOfBirth, emp.Address, emp.Experience, id)
	if err != nil {
		log.Printf("Error updating employee: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the employee has experience, update their previous company records
	if emp.Experience {
		// Fetch company records for the employee
		companies, err := prevcompanybyid(tx, id)
		if err != nil {
			log.Printf("Error retrieving company records: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update the first company record
		if len(companies) > 0 {
			_, err = tx.Exec(
				"UPDATE prevcompany SET companyName=?, position=?, startDate=?, endDate=?, duration=? WHERE id=?",
				emp.CompanyName,
				emp.Designation,
				emp.JoinDate,
				emp.RelievedDate,
				emp.TotalDuration,
				companies[0].id, // Use the first company's ID
			)
			if err != nil {
				log.Printf("Error updating first company: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Update the second company record if necessary
		if len(companies) > 1 && emp.SecondCompanyFormValue.CompanyName != "" {
			_, err = tx.Exec(
				"UPDATE prevcompany SET companyName=?, position=?, startDate=?, endDate=?, duration=? WHERE id=?",
				emp.SecondCompanyFormValue.CompanyName,
				emp.SecondCompanyFormValue.Designation,
				emp.SecondCompanyFormValue.JoinDate,
				emp.SecondCompanyFormValue.RelievedDate,
				emp.SecondCompanyFormValue.Duration,
				companies[1].id, // Use the second company's ID
			)
			if err != nil {
				log.Printf("Error updating second company: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {

		_, err = db.Exec("DELETE FROM prevcompany WHERE empId = ?", id)
		if err != nil {
			log.Printf("Error deleting from prevcompany: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	log.Printf("Employee updated successfully: %+v", emp)

	if emp.Experience {
		_, err = tx.Exec(
			"INSERT INTO prevcompany (companyName, position, startDate, endDate, duration, empId) VALUES (?, ?, ?, ?, ?, ?)",
			emp.CompanyName, emp.Designation, emp.JoinDate, emp.RelievedDate, emp.TotalDuration, id)
		if err != nil {
			log.Printf("Error inserting company record: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if emp.SecondCompanyFormValue.CompanyName != "" {
			_, err = tx.Exec(
				"INSERT INTO prevcompany (companyName, position, startDate, endDate, duration, empId) VALUES (?, ?, ?, ?, ?, ?)",
				emp.SecondCompanyFormValue.CompanyName, emp.SecondCompanyFormValue.Designation, emp.SecondCompanyFormValue.JoinDate, emp.SecondCompanyFormValue.RelievedDate, emp.SecondCompanyFormValue.Duration, id)
			if err != nil {
				log.Printf("Error inserting second company record: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	log.Printf("Employee updated successfully: %+v", emp)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Employee added successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

type Role struct {
	RoleID   int    `json:"role_id"`
	RoleName string `json:"roleName"`
}

func getRoleById(w http.ResponseWriter, r *http.Request) {
	// Extract role ID from URL parameters
	vars := mux.Vars(r)
	roleid, ok := vars["id"]
	if !ok {
		http.Error(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	// Open database connection
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Prepare the query
	query := `SELECT role_id,roleName FROM role WHERE role_id = ?`

	// Execute the query with parameter
	row := db.QueryRow(query, roleid)

	// Scan the result into a Role struct
	var role Role
	if err := row.Scan(&role.RoleID, &role.RoleName); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No role found with the given ID", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Send the result as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getRolesByDepartment handles requests to fetch roles by department ID
func getRolesByDepartment(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	deptID, ok := vars["id"]
	if !ok {
		http.Error(w, "Department ID is required", http.StatusBadRequest)
		return
	}

	// Query to get roles based on department ID
	query := `SELECT role_id, roleName FROM role WHERE dept_id = ?`
	rows, err := db.Query(query, deptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.RoleID, &role.RoleName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		roles = append(roles, role)
	}

	// Convert the result to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(roles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Department represents the department structure
type Department struct {
	DeptID         int    `json:"dept_id"`
	DepartmentName string `json:"department"`
}

// getDepartments handles requests to fetch all departments
func getDepartments(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	query := `SELECT dept_id, department FROM department`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var dept Department
		if err := rows.Scan(&dept.DeptID, &dept.DepartmentName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		departments = append(departments, dept)
	}
	// Convert the result to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(departments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getDepartmentsById(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	deptID, ok := vars["id"]
	if !ok {
		http.Error(w, "Department ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT dept_id,department FROM department where dept_id = ?`
	rows, err := db.Query(query, deptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var dept Department
		if err := rows.Scan(&dept.DeptID, &dept.DepartmentName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		departments = append(departments, dept)
	}
	// 	// Convert the result to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(departments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Manager represents the manager structure
type Manager struct {
	ManagerID   int    `json:"manager_id"`
	ManagerName string `json:"managerName"`
}

// getManagers handles requests to fetch all managers
func getManagers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `SELECT manager_id, managerName FROM manager`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var managers []Manager
	for rows.Next() {
		var mgr Manager
		if err := rows.Scan(&mgr.ManagerID, &mgr.ManagerName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		managers = append(managers, mgr)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(managers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getManagerByDepartment handles requests to fetch the manager by department ID
func getManagerByDepartment(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	deptID, ok := vars["id"]
	if !ok {
		http.Error(w, "Department ID is required", http.StatusBadRequest)
		return
	}
	// deptID := r.URL.Query().Get("dept_id")
	query := `SELECT m.manager_id, m.managerName 
              FROM manager m
              JOIN department d ON m.manager_id = d.manager_id
              WHERE d.dept_id = ?`
	row := db.QueryRow(query, deptID)

	var mgr Manager
	if err := row.Scan(&mgr.ManagerID, &mgr.ManagerName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mgr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getEmployeeAsManager(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	deptID, ok := vars["id"]
	if !ok {
		http.Error(w, "Department ID is required", http.StatusBadRequest)
		return
	}
	query := `select firstName,lastName from employeeRole where tech_lead is NUll and dept_id = ?`
	row := db.QueryRow(query, deptID)

	var mgr Employee
	if err := row.Scan(&mgr.FirstName, &mgr.LastName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mgr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type AssignEmployee struct {
	EmpID    int `json:"employee_id"`
	RoleID   int `json:"role_id"`
	DeptID   int `json:"dept_id"`
	TechLead int `json:"tech_lead"`
	// TechLead struct {
	// 	Type string `json:"type"`
	// 	ID   int    `json:"id"`
	// } `json:"tech_lead"`
}

func assignEmployee(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var ass AssignEmployee
	log.Printf("%v", ass)
	// Decode JSON request body into AssignEmployee struct
	err = json.NewDecoder(r.Body).Decode(&ass)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Assigning Employee: %+v", ass)

	// Check if the employee is already assigned to the same role and department
	var exists bool
	query := `SELECT EXISTS(SELECT 6 FROM employeeRole WHERE employee_id = ?)`
	err = db.QueryRow(query, ass.EmpID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking existing assignment: %v", err)
		http.Error(w, "Error checking existing assignment", http.StatusInternalServerError)
		return
	}

	if exists {
		log.Printf("Employee is already assigned to this role and department")
		http.Error(w, "Employee is already assigned to this role and department", http.StatusConflict)
		return
	}
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO employeeRole (employee_id, role_id, dept_id,tech_lead) VALUES (? , ? , ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Error preparing statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(ass.EmpID, ass.RoleID, ass.DeptID, ass.TechLead)
	if err != nil {
		log.Printf("Error executing insert statement: %v", err)
		http.Error(w, "Error executing insert statement", http.StatusInternalServerError)
		return
	}

	// Prepare the SQL statement for UPDATE
	stmt, err = db.Prepare("UPDATE emply SET role_id = ?, dept_id = ?, tech_lead = ? WHERE empId = ?")
	if err != nil {
		log.Printf("Error preparing update statement: %v", err)
		http.Error(w, "Error preparing update statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement for UPDATE
	_, err = stmt.Exec(ass.RoleID, ass.DeptID, ass.TechLead, ass.EmpID)
	if err != nil {
		log.Printf("Error executing update statement: %v", err)
		http.Error(w, "Error executing update statement", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ass); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type TechLead struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

type HierarchyData struct {
	EmpID       int    `json:"empId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	RoleName    string `json:"roleName"`
	RoleID      int    `json:"role_id"`
	Department  string `json:"department"`
	ManagerName string `json:"managerName"`
	TechLead    int    `json:"techLead"`
}

// Fetch hierarchy data from the database
func getHierarchyData(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// SQL query to fetch employee hierarchy data
	query := `
		SELECT 
			e.empId,
			e.firstName,
			e.lastName,
			r.roleName,
			r.role_id,
			d.department,
			m.managerName,
			er.tech_lead
		FROM 
			employeeRole er
		JOIN 
			emply e ON er.employee_id = e.empId
		JOIN 
			role r ON er.role_id = r.role_id
		JOIN 
			department d ON er.dept_id = d.dept_id
		JOIN 
			manager m ON d.manager_id = m.manager_id
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var hierarchyData []HierarchyData
	for rows.Next() {
		var techLeadJSON string
		var data HierarchyData

		if err := rows.Scan(&data.EmpID, &data.FirstName, &data.LastName, &data.RoleName, &data.RoleID, &data.Department, &data.ManagerName, &techLeadJSON); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		// Log the tech_lead JSON string
		log.Printf("TechLead JSON string: %s", techLeadJSON)
		// Deserialize tech_lead JSON string into TechLead struct
		if err := json.Unmarshal([]byte(techLeadJSON), &data.TechLead); err != nil {
			log.Printf("Error unmarshalling tech_lead JSON: %v", err)
			http.Error(w, "Error processing tech_lead data", http.StatusInternalServerError)
			return
		}

		hierarchyData = append(hierarchyData, data)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		http.Error(w, "Error reading data", http.StatusInternalServerError)
		return
	}

	// Send a success response with the hierarchy data
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hierarchyData); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func getHierarchyDataById(w http.ResponseWriter, r *http.Request) {
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["id"]

	// Convert the employee ID to an integer
	empId, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// SQL query to fetch employee hierarchy data by ID
	query := `
		SELECT 
			e.empId,
			e.firstName,
			e.lastName,
			r.roleName,
            r.role_id,
			d.department,
			m.managerName
		FROM 
			employeeRole er
		JOIN 
			emply e ON er.employee_id = e.empId
		JOIN 
			role r ON er.role_id = r.role_id
		JOIN 
			department d ON er.dept_id = d.dept_id
		JOIN 
			manager m ON d.manager_id = m.manager_id
		WHERE 
			er.role_id = ?
	`

	// Execute the query with the employee ID
	row := db.QueryRow(query, empId)

	var data HierarchyData
	if err := row.Scan(&data.EmpID, &data.FirstName, &data.LastName, &data.RoleName, &data.RoleID, &data.Department, &data.ManagerName); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Employee not found", http.StatusNotFound)
		} else {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
		}
		return
	}

	// Send a success response with the hierarchy data
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

type Documents struct {
	ID       int    `json:"id"`
	EmpID    int    `json:"empId"`
	Filename string `json:"filename"`
	FileData []byte `json:"filedata"`
}

// /////
func uploadDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	empId := r.FormValue("empId")

	// Extract files
	aadharFile, _, err := r.FormFile("aadhar")
	if err != nil {
		http.Error(w, "Unable to get aadhar file", http.StatusBadRequest)
		return
	}
	defer aadharFile.Close()

	profilePhoto, _, err := r.FormFile("profilephoto")
	if err != nil {
		http.Error(w, "Unable to get profile photo", http.StatusBadRequest)
		return
	}
	defer profilePhoto.Close()

	// Read files into memory
	aadharFileBytes, err := io.ReadAll(aadharFile)
	if err != nil {
		http.Error(w, "Unable to read aadhar file", http.StatusInternalServerError)
		return
	}

	profilePhotoBytes, err := io.ReadAll(profilePhoto)
	if err != nil {
		http.Error(w, "Unable to read profile photo", http.StatusInternalServerError)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert file data into database
	_, err = db.Exec("INSERT INTO documents (empId, aadharFile, profilePhoto) VALUES (?, ?, ?)",
		empId, aadharFileBytes, profilePhotoBytes)
	if err != nil {
		http.Error(w, "Error saving files to database", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Files uploaded successfully"))
}

type PersonalDetails struct {
	EmpID        int    `json:"empId"`
	Gender       string `json:"gender"`
	Relationship string `json:"relationship"`
	BloodGroup   string `json:"bloodgroup"`
}

func personaldetails(w http.ResponseWriter, r *http.Request) {
	var emp PersonalDetails

	// Decode the JSON request body into the PersonalDetails struct
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Printf("Error decoding JSON request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback if any error occurs before commit

	// Prepare the SQL insert statement
	stmt, err := tx.Prepare("INSERT INTO personaldetails (empId, gender, relationship) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the insert statement
	_, err = stmt.Exec(emp.EmpID, emp.Gender, emp.Relationship)
	if err != nil {
		log.Printf("Error executing insert statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Personal details added successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

// ////////////
func handlePersonalDetailsAndDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	empId := r.FormValue("empId")

	// Extract JSON data from form data
	jsonData := r.FormValue("personalDetails")
	if jsonData == "" {
		http.Error(w, "Missing personal details", http.StatusBadRequest)
		return
	}

	// Parse JSON data into PersonalDetails struct
	var emp PersonalDetails
	err = json.Unmarshal([]byte(jsonData), &emp)
	if err != nil {
		http.Error(w, "Error decoding JSON request body", http.StatusBadRequest)
		return
	}

	// Extract files
	aadharFile, _, err := r.FormFile("aadhar")
	if err != nil {
		http.Error(w, "Unable to get aadhar file", http.StatusBadRequest)
		return
	}
	defer aadharFile.Close()

	profilePhoto, _, err := r.FormFile("profilephoto")
	if err != nil {
		http.Error(w, "Unable to get profile photo", http.StatusBadRequest)
		return
	}
	defer profilePhoto.Close()

	// Read files into memory
	aadharFileBytes, err := io.ReadAll(aadharFile)
	if err != nil {
		http.Error(w, "Unable to read aadhar file", http.StatusInternalServerError)
		return
	}

	profilePhotoBytes, err := io.ReadAll(profilePhoto)
	if err != nil {
		http.Error(w, "Unable to read profile photo", http.StatusInternalServerError)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback if any error occurs before commit

	// Prepare and execute the SQL insert statement for personal details
	stmt, err := tx.Prepare("INSERT INTO personaldetails (empId, gender, relationship, bloodgroup) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(emp.EmpID, emp.Gender, emp.Relationship, emp.BloodGroup)
	if err != nil {
		log.Printf("Error executing personal details insert statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare and execute the SQL insert statement for file data
	_, err = tx.Exec("INSERT INTO documents (empId, aadharFile, profilePhoto) VALUES (?, ?, ?)",
		empId, aadharFileBytes, profilePhotoBytes)
	if err != nil {
		log.Printf("Error executing file data insert statement: %v", err)
		http.Error(w, "Error saving files to database", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Personal details and files added successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

func handleUpdatePersonalDetailsAndDocuments(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	empId := r.FormValue("empId")
	if empId == "" {
		http.Error(w, "Employee ID is missing", http.StatusBadRequest)
		return
	}

	// Extract JSON data from form data
	jsonData := r.FormValue("personalDetails")
	if jsonData == "" {
		http.Error(w, "Missing personal details", http.StatusBadRequest)
		return
	}

	// Parse JSON data into PersonalDetails struct
	var emp PersonalDetails
	err = json.Unmarshal([]byte(jsonData), &emp)
	if err != nil {
		http.Error(w, "Error decoding JSON request body", http.StatusBadRequest)
		return
	}

	// Extract files
	aadharFile, _, err := r.FormFile("aadhar")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Unable to get aadhar file", http.StatusBadRequest)
		return
	}
	defer aadharFile.Close()

	profilePhoto, _, err := r.FormFile("profilephoto")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Unable to get profile photo", http.StatusBadRequest)
		return
	}
	defer profilePhoto.Close()

	var aadharFileBytes []byte
	if aadharFile != nil {
		aadharFileBytes, err = io.ReadAll(aadharFile)
		if err != nil {
			http.Error(w, "Unable to read aadhar file", http.StatusInternalServerError)
			return
		}
	}

	var profilePhotoBytes []byte
	if profilePhoto != nil {
		profilePhotoBytes, err = io.ReadAll(profilePhoto)
		if err != nil {
			http.Error(w, "Unable to read profile photo", http.StatusInternalServerError)
			return
		}
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback if any error occurs before commit

	// Prepare and execute the SQL update statement for personal details
	stmt, err := tx.Prepare(`
        UPDATE personaldetails
        SET gender = ?, relationship = ?, bloodgroup = ?
        WHERE empId = ?`)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(emp.Gender, emp.Relationship, emp.BloodGroup, empId)
	if err != nil {
		log.Printf("Error executing personal details update statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare and execute the SQL update statement for file data
	_, err = tx.Exec(`
        UPDATE documents
        SET aadharFile = ?, profilePhoto = ?
        WHERE empId = ?`,
		aadharFileBytes, profilePhotoBytes, empId)
	if err != nil {
		log.Printf("Error executing file data update statement: %v", err)
		http.Error(w, "Error saving files to database", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Personal details and files updated successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

type DocumentResponse struct {
	AadharImage  string `json:"aadharImage"`
	ProfilePhoto string `json:"profilePhoto"`
	Gender       string `json:"gender"`
	Relationship string `json:"relationship"`
	BloodGroup   string `json:"bloodgroup"`
}

// type DocumentResponse struct {
// 	AadharImage  string `json:"aadharImage"`
// 	ProfilePhoto string `json:"profilePhoto"`
// }

// func getDocuments(w http.ResponseWriter, r *http.Request) {
// 	// Get the employee ID from the URL parameters
// 	vars := mux.Vars(r)
// 	empIdStr := vars["id"]

// 	// Convert the employee ID to an integer
// 	empId, err := strconv.Atoi(empIdStr)
// 	if err != nil {
// 		log.Printf("Invalid employee ID: %v", err)
// 		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Connect to the database
// 	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Fetch the document data from the database
// 	var aadharFileBytes, profilePhotoBytes []byte
// 	err = db.QueryRow("SELECT aadharFile, profilePhoto FROM documents WHERE empId = ?", empId).Scan(&aadharFileBytes, &profilePhotoBytes)
// 	if err != nil {
// 		http.Error(w, "Document not found", http.StatusNotFound)
// 		return
// 	}

// 	// Encode files to Base64
// 	aadharBase64 := base64.StdEncoding.EncodeToString(aadharFileBytes)
// 	profilePhotoBase64 := base64.StdEncoding.EncodeToString(profilePhotoBytes)

// 	// Create a JSON response
// 	response := DocumentResponse{
// 		AadharImage:  aadharBase64,
// 		ProfilePhoto: profilePhotoBase64,
// 	}

// 	// Set the Content-Type header to application/json
// 	w.Header().Set("Content-Type", "application/json")

// 	// Encode the response as JSON
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		log.Printf("Failed to encode JSON response: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// }

func getDocuments(w http.ResponseWriter, r *http.Request) {
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["id"]

	// Convert the employee ID to an integer
	empId, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Fetch the document and additional data from the database
	var aadharFileBytes, profilePhotoBytes []byte
	var gender, relationship, bloodGroup string
	err = db.QueryRow(`
		SELECT d.aadharFile, d.profilePhoto, a.gender,a.relationship,a.bloodgroup
		FROM documents d
		JOIN personaldetails a ON d.empId = a.empId
		WHERE d.empId = ?
	`, empId).Scan(&aadharFileBytes, &profilePhotoBytes, &gender, &relationship, &bloodGroup)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rows found for empId %d: %v", empId, err)
			http.Error(w, "Document or additional data not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to fetch data for empId %d: %v", empId, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Encode files to Base64
	aadharBase64 := base64.StdEncoding.EncodeToString(aadharFileBytes)
	profilePhotoBase64 := base64.StdEncoding.EncodeToString(profilePhotoBytes)

	// Create a JSON response
	response := DocumentResponse{
		AadharImage:  aadharBase64,
		ProfilePhoto: profilePhotoBase64,
		Gender:       gender,
		Relationship: relationship,
		BloodGroup:   bloodGroup,
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// /
// setManager handles fetching managers from the database based on department and tech lead
func setManager(w http.ResponseWriter, r *http.Request) {
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["id"]

	// Convert the employee ID to an integer
	id, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare the query with parameters
	query := `SELECT empId FROM emply WHERE dept_id=? AND tech_lead IS NULL`

	// Query the database with the provided department ID
	rows, err := db.Query(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var managers []Employee
	for rows.Next() {
		var mgr Employee
		if err := rows.Scan(&mgr.EmpId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		managers = append(managers, mgr)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(managers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getAssignData(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(` SELECT * FROM employeeRole`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var employees []AssignEmployee
	for rows.Next() {
		var emp AssignEmployee
		if err := rows.Scan(&emp.EmpID, &emp.RoleID, &emp.DeptID, &emp.TechLead); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved data: %v", employees)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}


/////******/////
////register employee
func registerHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert into database
	query := "INSERT INTO reg_employee (firstname, lastname, email, password) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, firstname, lastname, email, string(hashedPassword))
	if err != nil {
		log.Printf("Error saving user: %v", err)  // Detailed error logging
		http.Error(w, "Failed to register: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Registration successful")
}

var store = sessions.NewCookieStore([]byte("secret-key"))

func loginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    email := r.FormValue("email")
    password := r.FormValue("password")

    var hashedPassword string
    var empId int

    err = db.QueryRow("SELECT empId, password FROM reg_employee WHERE email = ?", email).Scan(&empId, &hashedPassword)
    if err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    session, _ := store.Get(r, "session-name")
    session.Values["empId"] = empId
    session.Save(r, w)

    // Track login time
    loginTime := time.Now()
    _, err = db.Exec("INSERT INTO employee_time_logs (empId, login_time) VALUES (?, ?)", empId, loginTime)
    if err != nil {
        http.Error(w, "Error tracking login", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Login successful")
}



func logoutHandler(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
    if err != nil {
        log.Printf("Error opening database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var emp Employee
    err = json.NewDecoder(r.Body).Decode(&emp)
    if err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Track logout time
    logoutTime := time.Now()

    var loginTimeStr string
    err = db.QueryRow("SELECT login_time FROM employee_time_logs WHERE empId = ? ORDER BY logId DESC LIMIT 1", emp.EmpId).Scan(&loginTimeStr)
    if err != nil {
        log.Printf("Error fetching login time: %v", err)
        http.Error(w, "Error fetching login time", http.StatusInternalServerError)
        return
    }

    // Parse the scanned login_time string into a time.Time object
    loginTime, err := time.Parse("2006-01-02 15:04:05", loginTimeStr)
    if err != nil {
        log.Printf("Error parsing login time: %v", err)
        http.Error(w, "Error parsing login time", http.StatusInternalServerError)
        return
    }

    workingHours := logoutTime.Sub(loginTime)

    // Convert workingHours (time.Duration) to a string format "HH:MM:SS"
    hours := int(workingHours.Hours())
    minutes := int(workingHours.Minutes()) % 60
    seconds := int(workingHours.Seconds()) % 60
    workingHoursStr := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

    _, err = db.Exec("UPDATE employee_time_logs SET logout_time = ?, working_hours = ? WHERE empId = ? ORDER BY logId DESC LIMIT 1", logoutTime, workingHoursStr, emp.EmpId)
    if err != nil {
        log.Printf("Error tracking logout: %v", err)
        http.Error(w, "Error tracking logout", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Logout successful, worked for %v", workingHoursStr)
}


	


func timesheetHandler(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/employee")
    if err != nil {
        log.Printf("Error opening database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["empId"]

	// Convert the employee ID to an integer
	empId, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}
    query := `
    SELECT e.empId, e.firstname, e.lastname, t.login_time, t.logout_time, t.working_hours
    FROM emply e
    JOIN employee_time_logs t ON e.empId = t.empId
    WHERE e.empId = ?`

    rows, err := db.Query(query, empId)
    if err != nil {
        log.Printf("Error fetching timesheet data: %v", err)
        http.Error(w, "Error fetching timesheet data", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var timesheet []map[string]interface{}

    for rows.Next() {
        var empId int
        var firstname, lastname, loginTime, logoutTime, workingHours string
        if err := rows.Scan(&empId, &firstname, &lastname, &loginTime, &logoutTime, &workingHours); err != nil {
            log.Printf("Error scanning timesheet data: %v", err)
            http.Error(w, "Error scanning timesheet data", http.StatusInternalServerError)
            return
        }

        timesheet = append(timesheet, map[string]interface{}{
            "empId":         empId,
            "firstname":     firstname,
            "lastname":      lastname,
            "login_time":    loginTime,
            "logout_time":   logoutTime,
            "working_hours": workingHours,
        })
    }

    if len(timesheet) == 0 {
        http.Error(w, "No timesheet data found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(timesheet)
}



////////

//main

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/employee", dataHandler).Methods("GET")
	r.HandleFunc("/addemployee", addEmployeeHandler).Methods("POST")
	// r.HandleFunc("/getemployee/employee/{id}", getEmployeeByID).Methods("GET")
	//
	// r.HandleFunc("/getbyid/prevcompany/{id}", prevcompanybyid).Methods("GET")
	r.HandleFunc("/getbyid/emply/{id}", emplybyid).Methods("GET")
	r.HandleFunc("/getbyid/emplywithprevcompany/{id}", employeeWithCompanyById).Methods("GET")
	//
	r.HandleFunc("/delete/employee/{id}", deleteEmployeeHandler).Methods("DELETE")
	r.HandleFunc("/update/employee/{id}", updateEmployeeHandler).Methods("PUT")
	r.HandleFunc("/getrolesbydepartment/{id}", getRolesByDepartment).Methods("GET")
	r.HandleFunc("/departments", getDepartments).Methods("GET")
	r.HandleFunc("/managers", getManagers).Methods("GET")
	r.HandleFunc("/manager/{id}", getManagerByDepartment).Methods("GET")
	r.HandleFunc("/assignemployee", assignEmployee).Methods("POST")
	r.HandleFunc("/uploaddocuments", uploadDocumentsHandler).Methods("POST")
	r.HandleFunc("/personaldetails", personaldetails).Methods("POST")
	r.HandleFunc("/handlePersonalDetailsAndDocuments", handlePersonalDetailsAndDocuments).Methods("POST")

	r.HandleFunc("/hierarchy", getHierarchyData).Methods("GET")
	// Register the route with the router
	r.HandleFunc("/hierarchy/{id}", getHierarchyDataById).Methods("GET")
	r.HandleFunc("/departmentbyid/{id}", getDepartmentsById).Methods("GET")
	r.HandleFunc("/getDocuments/{id}", getDocuments).Methods("GET")
	r.HandleFunc("/handleUpdatePersonalDetailsAndDocuments", handleUpdatePersonalDetailsAndDocuments).Methods("PUT")
	r.HandleFunc("/setManager/{id}", setManager).Methods("GET")
	r.HandleFunc("/getEmployeeAsManager/{id}", getEmployeeAsManager).Methods("GET")
	r.HandleFunc("/getRoleById/{id}", getRoleById).Methods("GET")
	r.HandleFunc("/getAssignData", getAssignData).Methods("GET")
	////****/////
	r.HandleFunc("/timesheet/{empId}", timesheetHandler).Methods("GET")
	r.HandleFunc("/logout", logoutHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/register", registerHandler).Methods("POST")


	

	

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
