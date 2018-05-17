-- Create Date : 2018-05-17T12:58:20 --
# Schema :  employees
## Tables
name | comments
:--- | :---
[departments](#departments) | 
[dept_emp](#deptemp) | 
[dept_manager](#deptmanager) | 
[employees](#employees) | 
[salaries](#salaries) | 
[titles](#titles) | 
<div style='page-break-after: always;'></div>


---
## Views
name |
:--- |
[current_dept_emp](#currentdeptemp) |
[dept_emp_latest_date](#deptemplatestdate) |
<div style='page-break-after: always;'></div>


---
## Functions
name | comments
:--- | :---
<div style='page-break-after: always;'></div>


---
## Procedures
name | comments
:--- | :---
<div style='page-break-after: always;'></div>


---
## departments
> Comment :   
> Engine : InnoDB  
> Collation : utf8_general_ci  

Columns
name | type | null | default | extra | comment
:--- | :--- | :--- | :--- | :--- | :---
dept_no | char(4) | NO |  |  | 
dept_name | varchar(40) | NO |  |  | 
Indexs
name | columns | isnull
:--- | :--- | :---
dept_name | dept_name | true 
PRIMARY | dept_no | true 
Create Script
```sql
CREATE TABLE `departments` (
  `dept_no` char(4) NOT NULL,
  `dept_name` varchar(40) NOT NULL,
  PRIMARY KEY (`dept_no`),
  UNIQUE KEY `dept_name` (`dept_name`)
) DEFAULT CHARSET=utf8
```
[goto table list...](#tables)
<div style='page-break-after: always;'></div>


## dept_emp
> Comment :   
> Engine : InnoDB  
> Collation : utf8_general_ci  

Columns
name | type | null | default | extra | comment
:--- | :--- | :--- | :--- | :--- | :---
emp_no | int(11) | NO |  |  | 
dept_no | char(4) | NO |  |  | 
from_date | date | NO |  |  | 
to_date | date | NO |  |  | 
Indexs
name | columns | isnull
:--- | :--- | :---
PRIMARY | emp_no,dept_no | true 
dept_no | dept_no | false 
Create Script
```sql
CREATE TABLE `dept_emp` (
  `emp_no` int(11) NOT NULL,
  `dept_no` char(4) NOT NULL,
  `from_date` date NOT NULL,
  `to_date` date NOT NULL,
  PRIMARY KEY (`emp_no`,`dept_no`),
  KEY `dept_no` (`dept_no`),
  CONSTRAINT `dept_emp_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE,
  CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
) DEFAULT CHARSET=utf8
```
[goto table list...](#tables)
<div style='page-break-after: always;'></div>


## dept_manager
> Comment :   
> Engine : InnoDB  
> Collation : utf8_general_ci  

Columns
name | type | null | default | extra | comment
:--- | :--- | :--- | :--- | :--- | :---
emp_no | int(11) | NO |  |  | 
dept_no | char(4) | NO |  |  | 
from_date | date | NO |  |  | 
to_date | date | NO |  |  | 
Indexs
name | columns | isnull
:--- | :--- | :---
PRIMARY | emp_no,dept_no | true 
dept_no | dept_no | false 
Create Script
```sql
CREATE TABLE `dept_manager` (
  `emp_no` int(11) NOT NULL,
  `dept_no` char(4) NOT NULL,
  `from_date` date NOT NULL,
  `to_date` date NOT NULL,
  PRIMARY KEY (`emp_no`,`dept_no`),
  KEY `dept_no` (`dept_no`),
  CONSTRAINT `dept_manager_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE,
  CONSTRAINT `dept_manager_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
) DEFAULT CHARSET=utf8
```
[goto table list...](#tables)
<div style='page-break-after: always;'></div>


## employees
> Comment :   
> Engine : InnoDB  
> Collation : utf8_general_ci  

Columns
name | type | null | default | extra | comment
:--- | :--- | :--- | :--- | :--- | :---
emp_no | int(11) | NO |  |  | 
birth_date | date | NO |  |  | 
first_name | varchar(14) | NO |  |  | 
last_name | varchar(16) | NO |  |  | 
gender | enum('M','F') | NO |  |  | 
hire_date | date | NO |  |  | 
Indexs
name | columns | isnull
:--- | :--- | :---
PRIMARY | emp_no | true 
Create Script
```sql
CREATE TABLE `employees` (
  `emp_no` int(11) NOT NULL,
  `birth_date` date NOT NULL,
  `first_name` varchar(14) NOT NULL,
  `last_name` varchar(16) NOT NULL,
  `gender` enum('M','F') NOT NULL,
  `hire_date` date NOT NULL,
  PRIMARY KEY (`emp_no`)
) DEFAULT CHARSET=utf8
```
[goto table list...](#tables)
<div style='page-break-after: always;'></div>


## salaries
> Comment :   
> Engine : InnoDB  
> Collation : utf8_general_ci  

Columns
name | type | null | default | extra | comment
:--- | :--- | :--- | :--- | :--- | :---
emp_no | int(11) | NO |  |  | 
salary | int(11) | NO |  |  | 
from_date | date | NO |  |  | 
to_date | date | NO |  |  | 
Indexs
name | columns | isnull
:--- | :--- | :---
PRIMARY | emp_no,from_date | true 
Create Script
```sql
CREATE TABLE `salaries` (
  `emp_no` int(11) NOT NULL,
  `salary` int(11) NOT NULL,
  `from_date` date NOT NULL,
  `to_date` date NOT NULL,
  PRIMARY KEY (`emp_no`,`from_date`),
  CONSTRAINT `salaries_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE
) DEFAULT CHARSET=utf8
```
[goto table list...](#tables)
<div style='page-break-after: always;'></div>


## titles
> Comment :   
> Engine : InnoDB  
> Collation : utf8_general_ci  

Columns
name | type | null | default | extra | comment
:--- | :--- | :--- | :--- | :--- | :---
emp_no | int(11) | NO |  |  | 
title | varchar(50) | NO |  |  | 
from_date | date | NO |  |  | 
to_date | date | YES |  |  | 
Indexs
name | columns | isnull
:--- | :--- | :---
PRIMARY | emp_no,title,from_date | true 
Create Script
```sql
CREATE TABLE `titles` (
  `emp_no` int(11) NOT NULL,
  `title` varchar(50) NOT NULL,
  `from_date` date NOT NULL,
  `to_date` date DEFAULT NULL,
  PRIMARY KEY (`emp_no`,`title`,`from_date`),
  CONSTRAINT `titles_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE
) DEFAULT CHARSET=utf8
```
[goto table list...](#tables)
<div style='page-break-after: always;'></div>


## current_dept_emp


Create Script
```sql
CREATE VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from (`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`) and (`l`.`to_date` = `d`.`to_date`))))
```
[goto view list...](#views)
<div style='page-break-after: always;'></div>


## dept_emp_latest_date


Create Script
```sql
CREATE VIEW `dept_emp_latest_date` AS select `dept_emp`.`emp_no` AS `emp_no`,max(`dept_emp`.`from_date`) AS `from_date`,max(`dept_emp`.`to_date`) AS `to_date` from `dept_emp` group by `dept_emp`.`emp_no`
```
[goto view list...](#views)
<div style='page-break-after: always;'></div>


