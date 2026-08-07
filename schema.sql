CREATE TABLE
    Players (
        PiD INTEGER NOT NULL PRIMARY KEY ,
        KiD INTEGER NOT NULL,
        dName TEXT,
        PFP TEXT,
        Alliance TEXT DEFAULT 'DES'
    );

CREATE TABLE
    Giftcode (
        code TEXT PRIMARY KEY,
        claimedAt TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );