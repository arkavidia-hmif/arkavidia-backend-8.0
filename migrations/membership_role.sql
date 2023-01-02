DO $$ BEGIN
    CREATE TYPE membership_role AS ENUM (
        'leader',
        'member'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$