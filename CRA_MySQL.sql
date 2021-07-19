CREATE database MYSTOREDB;
CREATE USER 'user2' @'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user2' @'localhost';
USE MYSTOREDB;
CREATE TABLE Users (
  UserName VARCHAR(30) NOT NULL PRIMARY KEY,
  Password VARCHAR(256),
  UUIDKey VARCHAR(256)
);
CREATE TABLE TypeVehicle (
  Type VARCHAR(50),
  Name VARCHAR(50) NOT NULL PRIMARY KEY,
  NumberOfSeats INT,
  Inventory INT,
  DailyCost FLOAT,
  WeeklyCost FLOAT,
  MonthlyCost FLOAT
);
INSERT INTO
  TypeVehicle (
    Type,
    Name,
    NumberOfSeats,
    Inventory,
    DailyCost,
    WeeklyCost,
    MonthlyCost
  )
VALUES
  ('Compact', 'Fit', 5, 0, 50, 320, 1250),
  ('Mini Vehicles', 'HUSTLER', 4, 0, 60, 380, 1500),
  ('Mini Vehicles', 'Move', 4, 0, 55, 350, 1400),
  ('MPV', 'Noah', 8, 0, 120, 750, 3000),
  ('Compact', 'Note', 5, 0, 55, 350, 1400),
  ('Hybrid', 'Note e-Power', 5, 0, 80, 500, 2000),
  ('Sedan', 'PREMIO', 5, 0, 45, 280, 1150),
  ('Sedan', 'Sienta', 6, 0, 40, 250, 1000),
  ('Mini Vehicles', 'Spacia', 4, 0, 55, 350, 1400),
  ('Compact', 'Vitz', 5, 0, 50, 320, 1250),
  ('MPV', 'VOXY', 8, 0, 120, 750, 3000),
  ('Compact', 'Yaris', 5, 0, 50, 320, 1250);
CREATE TABLE BookingsData (
    UserName VARCHAR(30),
    CarName VARCHAR(20) NOT NULL,
    StartDate DATE NOT NULL,
    EndDate DATE NOT NULL,
    Price FLOAT,
    DaysOfRenting INT,
    BookingID INT NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (BookingID)
  );