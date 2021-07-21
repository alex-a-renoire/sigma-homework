# A TCP client-server app for Persons

It lets clients write data to the database. The basic entity is a person with a certain name. 
Clients accept data in JSON format, e.g.:

- AddPerson {"func_name":"AddPerson", "data":{"name":"Bob"}}
- GetPerson {"func_name":"GetPerson", "data":{"id":1}}
- UpdatePerson {"func_name":"UpdatePerson", "data":{"id":1, "name":"Alice"}}
- DeletePerson {"func_name":"DeletePerson", "data":{"id":1}}

Cases accounted for: 

- Input is not a JSON
- Wrong data type of a json field: {"func_name":"GetPerson", "data":{"id":"1"}}
- Wrong json field tag name / absence of a required field: {"func":"GetPerson", "data":{"id":1}} / {"func_name":"GetPerson"} / {"data":{"id":0}}
- Wrong field value {"func_name":"wrong_func", "data":{"name":"Bob"}}
- Delete / get a person with non-existent ID