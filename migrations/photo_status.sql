DO $$ BEGIN
    CREATE TYPE photo_status AS ENUM (
        'waiting-for-approval',
        'approved',
        'denied'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$