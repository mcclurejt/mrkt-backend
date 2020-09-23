CREATE TABLE IF NOT EXISTS Candle
(
  id INT,
  date DATE NOT NULL,
  open FLOAT NOT NULL,
  high FLOAT NOT NULL,
  low FLOAT NOT NULL,
  close FLOAT NOT NULL,
  FOREIGN KEY (id) REFERENCES Ticker(id),
  PRIMARY KEY (id, date)
);

LOAD DATA LOCAL INFILE 'data/daily_adjusted/daily_adjusted_AAPL.csv'
  INTO TABLE Candle
  FIELDS TERMINATED BY ','
  LINES TERMINATED BY '\n'
  IGNORE 1 LINES
  (@date, open, high, low, @close, @adjustedclose, @volume, @dividend, @split)
  SET id = (SELECT id FROM Ticker WHERE name = "AAPL"), date = STR_TO_DATE(@date, '%Y-%m-%d'), close=@adjustedclose;

LOAD DATA LOCAL INFILE 'data/daily_adjusted/daily_adjusted_AMZN.csv'
  INTO TABLE Candle
  FIELDS TERMINATED BY ','
  LINES TERMINATED BY '\n'
  IGNORE 1 LINES
  (@date, open, high, low, @close, @adjustedclose, @volume, @dividend, @split)
  SET id = (SELECT id FROM Ticker WHERE name = "AMZN"), date = STR_TO_DATE(@date, '%Y-%m-%d'), close=@adjustedclose;

LOAD DATA LOCAL INFILE 'data/daily_adjusted/daily_adjusted_FB.csv'
  INTO TABLE Candle
  FIELDS TERMINATED BY ','
  LINES TERMINATED BY '\n'
  IGNORE 1 LINES
  (@date, open, high, low, @close, @adjustedclose, @volume, @dividend, @split)
  SET id = (SELECT id FROM Ticker WHERE name = "FB"), date = STR_TO_DATE(@date, '%Y-%m-%d'), close=@adjustedclose;

LOAD DATA LOCAL INFILE 'data/daily_adjusted/daily_adjusted_GOOG.csv'
  INTO TABLE Candle
  FIELDS TERMINATED BY ','
  LINES TERMINATED BY '\n'
  IGNORE 1 LINES
  (@date, open, high, low, @close, @adjustedclose, @volume, @dividend, @split)
  SET id = (SELECT id FROM Ticker WHERE name = "GOOG"), date = STR_TO_DATE(@date, '%Y-%m-%d'), close=@adjustedclose;

LOAD DATA LOCAL INFILE 'data/daily_adjusted/daily_adjusted_TSLA.csv'
  INTO TABLE Candle
  FIELDS TERMINATED BY ','
  LINES TERMINATED BY '\n'
  IGNORE 1 LINES
  (@date, open, high, low, @close, @adjustedclose, @volume, @dividend, @split)
  SET id = (SELECT id FROM Ticker WHERE name = "TSLA"), date = STR_TO_DATE(@date, '%Y-%m-%d'), close=@adjustedclose;
