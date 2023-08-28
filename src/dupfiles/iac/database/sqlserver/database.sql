-- Create the database
CREATE DATABASE nas
ON
(
    NAME = infra_data,
    FILENAME = '/data/nas_data.mdf',
    SIZE = 5000MB,
    MAXSIZE = 50000MB,
    FILEGROWTH = 1000MB
)
LOG ON
(
    NAME = infra_log,
    FILENAME = '/data/logs/nas_log.ldf',
    SIZE = 5000MB,
    MAXSIZE = 50000MB,
    FILEGROWTH = 1000MB
);

GO

-- Use the newly created database
USE nas;

GO

-- Create the FileInfo table
CREATE TABLE FileInfo (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Filename NVARCHAR(255),
    Location NVARCHAR(1024),
    Size BIGINT,
    CreatedDate DATETIME,
    ModifiedDate DATETIME,
    FileExtension NVARCHAR(50),
    HashValue NVARCHAR(64),
    Permissions NVARCHAR(10) NOT NULL,
    FileOwner NVARCHAR(64),
    [GUID] NVARCHAR(64) NOT NULL,
    IsDuplicate INT DEFAULT 0
);

GO