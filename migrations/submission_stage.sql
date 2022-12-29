DO $$ BEGIN
    CREATE TYPE submission_stage AS ENUM (
        'first-stage',
        'second-stage',
        'final-stage'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$