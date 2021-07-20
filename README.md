# A TCP client-server app for Persons

It lets clients write data to the database. The basic entity is a person with a certain name. 
Clients accept data in JSON format, e.g.:

`{"func_name":"AddPerson", "data":{"name":"Bob"}}`

Possible values for func_name:

- AddPerson
- UpdatePerson
- GetPerson
- DeletePerson

  