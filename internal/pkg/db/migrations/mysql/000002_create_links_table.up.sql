CREATE TABLE IF NOT EXISTS Links(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    Title VARCHAR(255),
    Address VARCHAR(127),
    UserID INT,
    FOREIGN KEY (UserID) REFERENCES Users(ID),
    PRIMARY KEY (ID)
)
