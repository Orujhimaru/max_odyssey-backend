-- Add html_table and svg_image columns to questions table
ALTER TABLE questions 
ADD COLUMN html_table TEXT,
ADD COLUMN svg_image TEXT;
