registry
---

registry for storing and providing protobuf files.

This is a thumbnail sketch of what needs to be done. Currently only Save and Get file functions are implemented, and it is up to the calling service to map those to stan channels, etc. Ultimately this service should provide advanced mappings such that request can be sent based on the channel and the correct proto files will be resolved and returned. Imports should also be mapped and handled correctly.

Currently the domain is reqiured to be provided when saving, but it could be possible to use the `package` part of the proto spec to handle that automatically.