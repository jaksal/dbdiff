dbdiff
======

create mysql schema and data diff script

install :
-------------

* install golang latest. (https://golang.org/dl/)
* download src (git clone https://github.com/jaksal/dbdiff.git)
* goto src folder 
* go get && go build

usage :
-------------

* create schema diff script

```
dbdiff -diff_type=schema
  -source="uid:pwd@tcp(server_ip:port)/dbname"
  -target="uid:pwd@tcp(server_ip:port)/dbname"  
  -output=output.sql
```

* create data diff scirpt

```
dbdiff -diff_type=data
  -source="uid:pwd@tcp(server_ip:port)/dbname" 
  -target="uid:pwd@tcp(server_ip:port)/dbname" 
  -output=output.sql
```

* create markdown document 

```
dbdiff -diff_type=doc
  -source="uid:pwd@tcp(server_ip:port)/dbname"
  -output=output.md
```

* extra option

```
  -include="test_"    // include db object name containing test_  
  -exclude="test_"    // exclude db object name containing test_  
  -output="[DATE]/output.sql"   // create today date folder
```

* bugs

not detect column rename in schema diff mode. make drop and add column script 
