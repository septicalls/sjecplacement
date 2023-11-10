# Routes for the web application

This file lists all planned routes that the web application will act on. It is split into two phases for development.

## Phase 1
- `/` Home Page. Lists all placement drives (upcoming ones too) in the order of most recent to least recent. Pagination possible. Upcoming placement drives highlighted.
- `/drive/{drive-id}` Lists details of the specific placement drive.
- `/create` GET Shows a form to add a new placement drive to DB.
- `/create` POST Validates and adds a new placement drive to DB.

## Phase 2 (Subject to change)
- `/` and `/drive/{drive_id}` routes are restricted to users who have the student role (logged in).
- `/create` GET and POST restricted to users who have the admin role.
- `/login` for authenticating users.
- `/logout`