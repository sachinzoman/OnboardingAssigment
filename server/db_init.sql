-- create -- Create a new table called 'restaurant' in schema 'restaurantDB'
-- Drop the table if it already exists
IF OBJECT_ID('restaurantDB.restaurant', 'U') IS NOT NULL
DROP TABLE restaurantDB.restaurant;
GO
-- Create the table in the specified schema
CREATE TABLE restaurant
(
    restaurantId INT NOT NULL AUTO_INCREMENT PRIMARY KEY , -- primary key column
    restaurantName VARCHAR(30) NOT NULL,
    rating FLOAT NOT NULL,
    cusines VARCHAR(50) NOT NULL,
    address VARCHAR(100) NOT NULL,
    startime TIME NOT NULL,
    endtime TIME NOT NULL,
    cft FLOAT NOT NULL,
    img_url VARCHAR(100)
);
-- GO

INSERT INTO restaurant
(restaurantName, rating, cusines, address, startime, endtime, cft, img_url)
VALUES
("Rajinder Da Dhaba", 4.1, "North Indian, Rolls", 
"AB-14B, Nauroji Nagar Marg, Opposite, Safdarjung Enclave, New Delhi, Delhi 110029", 
"17:24:00","00:30:00",800,
"https://b.zmtcdn.com/data/pictures/9/7319/e1b7673ed0aa2993b55b177409d5596c.jpg");

