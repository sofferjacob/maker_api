# Events

Events can be sent to the API to then generate analytics info. Events are stored in the events table, using only the columns they need (id, event_type and timestamp are required). Additional event data may be passed in the body field as a json.
The following events can be sent:

### user_register

This event will be sent automatically by the server when the `/id/register` endpoint is used. Fields:

No additional fields

### user_login

This event will be sent automatically by the server when the `/id/login` endpoint is used. Fields:

- uid

### draft_create

This event will be sent automatically by the server when a draft creation endpoint is used. Fields:

- uid
- level_id (may be null)
- draft_id

### draft_update

This event will be sent automatically by the server when a draft update endpoint is used. Fields:

- uid
- draft_id

### draft_delete

This event will be sent automatically by the server when a draft deletion endpoint is used. Fields:

- uid
- draft_id (may be null)

### level_create

This event will be sent automatically by the server when a level creation endpoint is used. Fields:

- uid
- level_id

### level_update

This event will be sent automatically by the server when a level update endpoint is used. Fields:

- uid
- level_id

### level_delete

This event will be sent automatically by the server when a level deletion endpoint is used. Fields:

- uid
- level_id

### level_clone

This event will be sent automatically by the server when a level is cloned through `/drafts/level/:id`. Fields:

- uid (cloner)
- level_id

### game_start

This event should be sent when a user starts playing a level. Fields:

- uid
- level_id

### game_finish

This event should be sent upon gameplay completion. Fields:

- uid
- level_id
- state
- time [> 0] (time it took the player to finish the level, in seconds).
