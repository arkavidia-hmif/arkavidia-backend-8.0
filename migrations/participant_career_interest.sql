DO $$ BEGIN
    CREATE TYPE participant_career_interest AS ENUM (
        'software-engineering',
        'product-management',
        'ui-designer',
        'ux-designer',
        'ux-researcher',
        'it-consultant',
        'game-developer',
        'cyber-security',
        'business-analyst',
        'business-intelligence',
        'data-scientist',
        'data-analyst'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$