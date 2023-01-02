DO $$ BEGIN
    CREATE TYPE team_status AS ENUM (
        'waiting-for-evaluation',
        'passed',
        'eliminated'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$