dbdiff
======

create mysql schema and data diff script.
create markdown schema document.

mysql database 를 비교하여 스키마 및 데이터 변경사항을 sql 로 추출해주는 툴입니다.

스키마 문서를 md 파일로 생성하는 기능도 있습니다.

golang으로 작성되어 어느 OS에서도 빌드해서 사용가능합니다.

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
-- 샘플파일은 sample 폴더에 있습니다. 
-- output md file convert to pdf ==> ( https://github.com/jaksal/md2pdf )

* extra option

```
  -include="test_"    // include db object name containing test_  
  -exclude="test_"    // exclude db object name containing test_  
  -output="[DATE]/output.sql"   // create today date folder
  -ignore_column=update_date  // ignore column in data diff mode
```

* bug

테이블 비교시 컬럼이름만 변경된 경우는 sql 만으로 알수가 없어서 해당컬럼을 drop 하고 새로 추가합니다. 
이 과정에서 데이터가 유실되니 이부분 유의해주시기 바랍니다.
