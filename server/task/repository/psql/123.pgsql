SELECT * FROM tasks;
SELECT * FROM users;

DELETE FROM tasks;
DELETE FROM users;

UPDATE public.tasks
    SET title = 'new title',
        update_date = CURRENT_TIMESTAMP
    WHERE id = 65;

INSERT INTO users ("login", "hashpass") VALUES ('test', '1234');

SELECT CURRENT_TIMESTAMP;