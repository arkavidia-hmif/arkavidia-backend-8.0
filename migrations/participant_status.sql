DO $$ BEGIN
    CREATE TYPE participant_status AS ENUM (
        'waiting-for-verification',
        'verified',
        'declined'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$