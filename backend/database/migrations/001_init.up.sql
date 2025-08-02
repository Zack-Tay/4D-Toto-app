CREATE TABLE draw_results (
    id SERIAL PRIMARY KEY,
    draw_type VARCHAR(10) NOT NULL CHECK (draw_type IN ('4D', 'TOTO')),
    draw_date DATE NOT NULL,
    draw_number INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

CREATE UNIQUE INDEX idx_unqiue_draw ON draw_results(draw_type, draw_date, draw_number);
CREATE INDEX idx_draw_date ON draw_results(draw_date DESC);

CREATE TABLE results_4d (
    id SERIAL PRIMARY KEY,
    draw_results_id INTEGER NOT NULL REFERENCES draw_results(id) ON DELETE CASCADE,
    prize_category VARCHAR(20) NOT NULL, -- for example: 1st, 2nd or starters etc.
    position INTEGER NOT NULL DEFAULT 1,
    winning_number VARCHAR(4) NOT NULL
)

CREATE INDEX idx_4d_draw_date on results_4d(draw_results_id);

CREATE TABLE results_toto (
    id SERIAL PRIMARY KEY,
    draw_results_id INTEGER NOT NULL REFERENCES draw_results(id) ON DELETE CASCADE,
    winning_numbers TEXT NOT NULL, -- comma separated version of 6 numbers.
    additional_number INTEGER NOT NULL
)

CREATE INDEX idx_toto_draw_date on results_toto(draw_results_id);