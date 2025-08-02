CREATE DATABASE sg_lottery;

-- Connect to the database and run:
-- \c sg_lottery
-- \i migrations/001_init.up.sql

-- Insert sample 4D data
INSERT INTO draw_results(draw_type, draw_date, draw_number)
VALUES ('4D', '2024-01-15', 4890);

-- Get the ID of the draw we just inserted
-- Insert 4D prizes (using the draw_result_id from above)
INSERT INTO results_4d (draw_result_id, prize_category, position, winning_number) VALUES
(1, '1st', 1, '1234'),
(1, '2nd', 1, '5678'),
(1, '3rd', 1, '9012'),
(1, 'Starter', 1, '1111'),
(1, 'Starter', 2, '2222'),
(1, 'Starter', 3, '3333'),
(1, 'Starter', 4, '4444'),
(1, 'Starter', 5, '5555'),
(1, 'Starter', 6, '6666'),
(1, 'Starter', 7, '7777'),
(1, 'Starter', 8, '8888'),
(1, 'Starter', 9, '9999'),
(1, 'Starter', 10, '0000'),
(1, 'Consolation', 1, '1001'),
(1, 'Consolation', 2, '2002'),
(1, 'Consolation', 3, '3003'),
(1, 'Consolation', 4, '4004'),
(1, 'Consolation', 5, '5005'),
(1, 'Consolation', 6, '6006'),
(1, 'Consolation', 7, '7007'),
(1, 'Consolation', 8, '8008'),
(1, 'Consolation', 9, '9009'),
(1, 'Consolation', 10, '0010');

-- Insert sample TOTO data
INSERT INTO draw_results (draw_type, draw_date, draw_number) 
VALUES ('TOTO', '2024-01-15', 3890);

INSERT INTO results_toto (draw_result_id, winning_numbers, additional_number) 
VALUES (2, '3,15,23,28,35,42', 7);