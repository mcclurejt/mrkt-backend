CREATE DATABASE IF NOT EXISTS ticker_data;

USE ticker_data;

SET GLOBAL local_infile = 'ON';

CREATE TABLE IF NOT EXISTS Ticker
(
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR (8) NOT NULL UNIQUE
);

INSERT INTO Ticker (name) VALUES ("AAPL");
INSERT INTO Ticker (name) VALUES ("AMZN");
INSERT INTO Ticker (name) VALUES ("FB");
INSERT INTO Ticker (name) VALUES ("GOOG");
INSERT INTO Ticker (name) VALUES ("TSLA");
INSERT INTO Ticker (name) VALUES ("BTC");

SOURCE data/daily_adjusted/insert_candles.sql;
