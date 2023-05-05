CREATE TABLE Employee
    (EmpID INT NOT NULL , 
        EmpName VARCHAR(50) NOT NULL, 
        Designation VARCHAR(50) NULL, 
        Department VARCHAR(50) NULL, 
        JoiningDate DATE NULL
    );
    

-- ADD ROWS TO THE TABLE.
    -- SQL SERVER 2008 AND ABOVE.

INSERT INTO Employee
    (EmpID, EmpName, Designation, Department, JoiningDate)
VALUES 
    (1, 'CHIN YEN', 'LAB ASSISTANT', 'LAB', CURRENT_DATE),
    (2, 'MIKE PEARL', 'SENIOR ACCOUNTANT', 'ACCOUNTS', CURRENT_DATE),
    (3, 'GREEN FIELD', 'ACCOUNTANT', 'ACCOUNTS', CURRENT_DATE),
    (4, 'DEWANE PAUL', 'PROGRAMMER', 'IT', CURRENT_DATE),
    (5, 'MATTS', 'SR. PROGRAMMER', 'IT', CURRENT_DATE),
    (6, 'PLANK OTO', 'ACCOUNTANT', 'ACCOUNTS', CURRENT_DATE);


-- SQL SERVER 2005 AND BEFORE.

INSERT INTO Employee (EmpID, EmpName, Designation, Department, JoiningDate)
    SELECT 1, 'CHIN YEN', 'LAB ASSISTANT', 'LAB', CURRENT_DATE;
INSERT INTO Employee (EmpID, EmpName, Designation, Department, JoiningDate)
    SELECT 2, 'MIKE PEARL', 'SENIOR ACCOUNTANT', 'ACCOUNTS', CURRENT_DATE;
INSERT INTO Employee (EmpID, EmpName, Designation, Department, JoiningDate)
    SELECT 3, 'GREEN FIELD', 'ACCOUNTANT', 'ACCOUNTS', CURRENT_DATE;
INSERT INTO Employee (EmpID, EmpName, Designation, Department, JoiningDate)
    SELECT 4, 'DEWANE PAUL', 'PROGRAMMER', 'IT', CURRENT_DATE;
INSERT INTO Employee (EmpID, EmpName, Designation, Department, JoiningDate)
    SELECT 5, 'MATTS', 'SR. PROGRAMMER', 'IT', CURRENT_DATE;
INSERT INTO Employee (EmpID, EmpName, Designation, Department, JoiningDate)
    SELECT 6, 'PLANK OTO', 'ACCOUNTANT', 'ACCOUNTS', CURRENT_DATE;
    
