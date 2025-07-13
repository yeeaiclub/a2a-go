# changelog

## v0.2.1
1. Fix issue where task ID is empty and not properly set in message parameters
2. Fix issue where TaskArtifactUpdateEvent is not properly included in task processing


## v0.2.0
1. Fix json: cannot unmarshal object into Go struct field Message.message.parts of type types.Part (#5)
2. Add APIKey, HTTP, OAuth2, and OpenIdConnect security schemes, structures and constants.(#7)
3. Refactor(server/event/queue): introduce eventType for state transitions. Removed consumer.go and its tests; (#8)
4. Add client auth middleware, update types and docs, add issue templates (#10)
5. Refactor: adapt context handling across server components and updated readme (#12)

