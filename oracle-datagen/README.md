# Oracle Data Generator
This project will contain some test-data generation scripts to pump Oracle databases with large data to test backups. This project will use Oracle SQL*Loader as a dependency. 

Currently this utility can populate the data to some pre-defined tables. User cannot choose random table name and random table fields to load the data here.
Below are the list of supported tables:
-calls
-users
-shipments
-bookings
-inventory
-ProductDescription
-payments
-orders
-products
-customers
-suppliers
-admin



Below is the respective schema that should be created before using this utility.
(NOTE:You can create these tables in any tablespace that you want, provided table names and fields are exactly the same as below.)
```
CREATE TABLE calls (
  call_id NUMBER,
  employee_id NUMBER,
  call_date VARCHAR(1000),
  call_duration NUMBER,
  customer_name VARCHAR2(100),
  customer_phone VARCHAR2(20),
  call_outcome CLOB, 
  customer_feedback CLOB
)
TABLESPACE USERS;

CREATE TABLE ProductDescription (
  name VARCHAR2(100) NOT NULL,
  description VARCHAR2(4000)
)
TABLESPACE SYSTEM;
	
CREATE TABLE bookings (
  user_id NUMBER NOT NULL,
  product_id NUMBER NOT NULL,
  quantity NUMBER NOT NULL,
  total NUMBER NOT NULL,
  order_date VARCHAR2(20) NOT NULL,
  description VARCHAR2(4000)
)
TABLESPACE products

CREATE TABLE admin (
  first_name VARCHAR2(50) NOT NULL,
  last_name VARCHAR2(50) NOT NULL,
  email VARCHAR2(100) NOT NULL,
  age NUMBER NOT NULL,
  description VARCHAR2(4000)
)
TABLESPACE Admin;

CREATE TABLE users (
  first_name VARCHAR2(50) NOT NULL,
  last_name VARCHAR2(50) NOT NULL,
  email VARCHAR2(100) NOT NULL,
  age NUMBER NOT NULL,
  description CLOB
)
TABLESPACE USERS;

CREATE TABLE products (
  name VARCHAR2(100) NOT NULL,
  description VARCHAR2(4000),
  price NUMBER NOT NULL,
  category VARCHAR2(50) NOT NULL
)
TABLESPACE products;

CREATE TABLE orders (
  user_id NUMBER NOT NULL,
  product_id NUMBER NOT NULL,
  quantity NUMBER NOT NULL,
  total NUMBER NOT NULL,
  order_date VARCHAR2(20) NOT NULL,
  description VARCHAR2(4000)
)
TABLESPACE products;

CREATE TABLE customers (
  first_name VARCHAR2(50) NOT NULL,
  last_name VARCHAR2(50) NOT NULL,
  email VARCHAR2(100) NOT NULL,
  phone VARCHAR2(20) NOT NULL,
  description VARCHAR2(4000)
)
TABLESPACE USERS;

CREATE TABLE suppliers (
  name VARCHAR2(100) NOT NULL,
  address VARCHAR2(200) NOT NULL,
  city VARCHAR2(50) NOT NULL,
  state VARCHAR2(50) NOT NULL,
  zip_code VARCHAR2(20) NOT NULL,
  country VARCHAR2(50) NOT NULL,
  phone VARCHAR2(20) NOT NULL,
  description VARCHAR2(4000)
)
TABLESPACE LogisticsSupplier;

CREATE TABLE inventory (
  product_id NUMBER NOT NULL,
  quantity NUMBER NOT NULL,
  reorder_lvl NUMBER NOT NULL,
  description VARCHAR2(4000)
)
TABLESPACE SYSTEM;

CREATE TABLE shipments (
  order_id NUMBER NOT NULL,
  tracking_no VARCHAR2(50) NOT NULL,
  carrier VARCHAR2(50) NOT NULL,
  ship_date VARCHAR2(50) NOT NULL,
  delivery_date VARCHAR2(20) NOT NULL,
  description VARCHAR2(4000)
)
TABLESPACE shipment;

CREATE TABLE payments (
  order_id NUMBER NOT NULL,
  amount NUMBER NOT NULL,
  payment_method VARCHAR2(50) NOT NULL,
  payment_date VARCHAR2(20) NOT NULL,
  description VARCHAR2(4000)
)
TABLESPACE shipment;
```