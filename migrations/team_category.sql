DO $$ BEGIN
    CREATE TYPE team_category AS ENUM (
        'competitive-programming',
        'datavidia',
        'uxvidia',
        'arkalogica'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$